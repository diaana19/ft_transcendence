package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/Transcendence/models"
	redispub "github.com/Transcendence/redis"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		allowed := []string{"http://localhost:3000", "http://localhost"}
		for _, a := range allowed {
			if origin == a {
				return true
			}
		}
		return false
	},
}

type ChatHandler struct {
	manager             *WSManager
	rdb                 *redis.Client
	notificationService *services.NotificationService
	fileRepo            repositories.FileRepository
	subscribedRooms     map[string]bool
	subscribedMu        sync.Mutex
}

type IncomingMessage struct {
	Action   string  `json:"action"`
	RoomID   string  `json:"room_id"`
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id"`
	FileID   *string `json:"file_id,omitempty"`
}

type OutgoingMessage struct {
	Type     string          `json:"type"`
	Message  *models.Message `json:"message,omitempty"`
	Username string          `json:"username,omitempty"`
	UserID   string          `json:"user_id,omitempty"`
	RoomID   string          `json:"room_id,omitempty"`
}

func NewChatHandler(
	manager *WSManager,
	rdb *redis.Client,
	notifService *services.NotificationService,
	fileRepo repositories.FileRepository,
) *ChatHandler {
	return &ChatHandler{
		manager:             manager,
		rdb:                 rdb,
		notificationService: notifService,
		fileRepo:            fileRepo,
		subscribedRooms:     make(map[string]bool),
	}
}


func (h *ChatHandler) sendPendingNotifications(client *Client) {
	notifs, err := h.notificationService.GetUnread(client.ID)
	if err != nil || len(notifs) == 0 {
		return
	}
	for _, n := range notifs {
		payload, err := json.Marshal(map[string]interface{}{
			"type":         "notification",
			"notification": n,
		})
		if err != nil {
			continue
		}
		safeSend(client.Send, payload)
	}
	h.notificationService.MarkAllRead(client.ID)
}

func (h *ChatHandler) HandleWS(c *gin.Context) {

	var userID string
	var username string
	if id, exists := c.Get("user_id"); exists {
		userID = id.(string)
		if u, ok := c.Get("username"); ok {
			username = u.(string)
		}
	} else {

		token := c.Query("token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}
		userID = claims.UserId
		username = claims.Username
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	client := &Client{
		ID:       userID,
		Username: username,
		Conn:     conn,
		Send:     make(chan []byte, 256),
	}

	h.manager.RegisterClient(client)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	defer h.manager.UnregisterClient(client)

	log.Printf("[WS] Client connected username=%q userID=%s , subscribing to notifications:%s", client.Username, client.ID, client.ID)
	redispub.Subscribe(ctx, h.rdb, "notifications:"+client.ID, func(payload string) {
		log.Printf("[WS] Forwarding notification to client username=%q userID=%s ", client.Username, client.ID)
		safeSend(client.Send, []byte(payload))
	})
	h.sendPendingNotifications(client)

	go client.WritePump()

	h.readPump(client)
}

func (h *ChatHandler) readPump(client *Client) {
	defer client.Conn.Close()

	for {
		_, raw, err := client.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close for client %s: %v\n", client.ID, err)
			}
			break
		}
		h.HandleMessage(client, raw)
	}
}

func (h *ChatHandler) HandleMessage(client *Client, raw []byte) {
	var incoming IncomingMessage
	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Printf("Invalid message format from %s: %v\n", client.ID, err)
		return
	}

	switch incoming.Action {
	case "join":
		h.handleJoin(client, incoming.RoomID)
	case "leave":
		h.handleLeave(client, incoming.RoomID)
	case "message":
		h.handleChat(client, incoming)
	default:
		log.Printf("Unknown action from %s: %s\n", client.ID, incoming.Action)
	}
}

func (h *ChatHandler) handleJoin(client *Client, roomID string) {
	if roomID == "" {
		return
	}
	h.manager.JoinRoom(client, roomID)

	h.subscribedMu.Lock()
	if !h.subscribedRooms[roomID] {
		h.subscribedRooms[roomID] = true
		redispub.Subscribe(context.Background(), h.rdb, "chat:"+roomID, func(payload string) {
			h.manager.BroadcastToRoom(roomID, []byte(payload), "")
		})
	}
	h.subscribedMu.Unlock()

	out := OutgoingMessage{
		Type:     "joined",
		UserID:   client.ID,
		Username: client.Username,
		RoomID:   roomID,
	}
	h.publishToRoom(roomID, out)
}

func (h *ChatHandler) handleLeave(client *Client, roomID string) {
	if roomID == "" {
		return
	}
	h.manager.LeaveRoom(client, roomID)

	out := OutgoingMessage{
		Type:     "left",
		UserID:   client.ID,
		Username: client.Username,
		RoomID:   roomID,
	}
	h.publishToRoom(roomID, out)
}

func (h *ChatHandler) handleChat(client *Client, incoming IncomingMessage) {
	if incoming.Content == "" && incoming.FileID == nil {
		return
	}
	if incoming.RoomID == "" {
		return
	}

	if incoming.FileID != nil {
		if err := h.handleAttachment(client.ID, incoming.RoomID, *incoming.FileID); err != nil {
			log.Printf("[Chat] attachment rejected: %v", err)
			return
		}
	}

	msg := models.Message{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		SenderID:  client.ID,
		Username:  client.Username,
		RoomID:    incoming.RoomID,
		Content:   incoming.Content,
		ParentID:  incoming.ParentID,
		FileID:    incoming.FileID,
		Type:      "dm",
	}

	out := OutgoingMessage{
		Type:    "message",
		Message: &msg,
	}

	h.publishToRoom(incoming.RoomID, out)
}

func (h *ChatHandler) publishToRoom(roomID string, out OutgoingMessage) {
	payload, err := json.Marshal(out)
	if err != nil {
		log.Printf("Marshal error: %v\n", err)
		return
	}

	if err := redispub.Publish(h.rdb, "chat:"+roomID, string(payload)); err != nil {
		log.Printf("Publish error: %v\n", err)
	}
}

func (h *ChatHandler) handleAttachment(senderID, roomID, fileID string) error {
	file, err := h.fileRepo.GetByID(fileID)
	if err != nil {
		return fmt.Errorf("file not found: %w", err)
	}
	if file.OwnerID != senderID {
		return fmt.Errorf("file not owned by sender")
	}
	if file.Visibility != models.FileVisibilityPrivate {
		return fmt.Errorf("file must be uploaded with visibility=private for DM attachments")
	}

	parts := strings.Split(roomID, ":")
	if len(parts) != 3 || parts[0] != "dm" {
		return fmt.Errorf("invalid room id format for DM attachment (expected dm:userA:userB)")
	}

	var recipientID string
	if parts[1] == senderID {
		recipientID = parts[2]
	} else if parts[2] == senderID {
		recipientID = parts[1]
	} else {
		return fmt.Errorf("sender not part of this room")
	}

	if err := h.fileRepo.GrantAccess(fileID, recipientID); err != nil {
		return fmt.Errorf("failed to grant access: %w", err)
	}

	log.Printf("[Chat] attachment granted: fileID=%s, sender=%s → recipient=%s",
		fileID, senderID, recipientID)
	return nil
}
