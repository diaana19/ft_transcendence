package socket

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/models"
	redispub "ft_transcendence/backend/internal/redis"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
)

type ChatHandler struct {
	manager             *WSManager
	rdb                 *redis.Client
	notificationService *services.NotificationService
	msgRepo             repositories.MessageRepository
	fileRepo            repositories.FileRepository
	upgrader            websocket.Upgrader
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

// NewChatHandler builds a ChatHandler whose WebSocket upgrader accepts
// connections from frontendURL (the SPA's public origin) and from clients
// that omit the Origin header (non-browser tools, smoke tests). "null" is
// also allowed for sandboxed iframes and file:// origins during local dev.
func NewChatHandler(
	manager *WSManager,
	rdb *redis.Client,
	notifService *services.NotificationService,
	msgRepo repositories.MessageRepository,
	fileRepo repositories.FileRepository,
	frontendURL string,
) *ChatHandler {
	allowed := []string{frontendURL, "null", ""}
	return &ChatHandler{
		manager:             manager,
		rdb:                 rdb,
		notificationService: notifService,
		msgRepo:             msgRepo,
		fileRepo:            fileRepo,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return slices.Contains(allowed, r.Header.Get("Origin"))
			},
		},
		subscribedRooms: make(map[string]bool),
	}
}

func (h *ChatHandler) sendPendingNotifications(client *Client) {
	notifs, err := h.notificationService.GetUnread(client.ID)
	if err != nil || len(notifs) == 0 {
		return
	}
	allDelivered := true
	for _, n := range notifs {
		payload, err := json.Marshal(map[string]any{
			"type":         "notification",
			"notification": n,
		})
		if err != nil {
			allDelivered = false
			continue
		}
		if !safeSend(client.Send, payload) {
			allDelivered = false
		}
	}
	// Only flag the backlog as read if every notification actually reached the
	// client's send buffer. If anything dropped, leave them unread so they
	// resurface on next connect rather than disappearing silently.
	if allDelivered {
		_ = h.notificationService.MarkAllRead(client.ID)
	}
}

func (h *ChatHandler) HandleWS(c *gin.Context) {
	userID := c.GetString("user_id")
	username := c.GetString("username")

	conn, err := h.upgrader.Upgrade(c.Writer, c.Request, nil)
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
	// LIFO: cancel fires first so the notifications subscriber goroutine starts
	// winding down before UnregisterClient closes client.Send.
	defer h.manager.UnregisterClient(client)
	defer cancel()

	log.Printf("[WS] Client connected username=%q userID=%s , subscribing to notifications:%s",
		client.Username, client.ID, client.ID)
	redispub.Subscribe(ctx, h.rdb, "notifications:"+client.ID, func(payload string) {
		log.Printf("[WS] Forwarding notification to client username=%q userID=%s ", client.Username, client.ID)
		safeSend(client.Send, []byte(payload))
	})
	h.sendPendingNotifications(client)

	go client.WritePump()

	h.readPump(client)
}

func (h *ChatHandler) readPump(client *Client) {
	defer func() { _ = client.Conn.Close() }()

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

	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("uuid error: %v", err)
		return
	}

	msgType := "room"
	if strings.HasPrefix(incoming.RoomID, "dm:") {
		msgType = "dm"
	}

	msg := models.Message{
		ID:        id.String(),
		CreatedAt: time.Now(),
		SenderID:  client.ID,
		Username:  client.Username,
		RoomID:    incoming.RoomID,
		Content:   incoming.Content,
		ParentID:  incoming.ParentID,
		FileID:    incoming.FileID,
		Type:      msgType,
	}

	if err := h.msgRepo.Create(&msg); err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	if msgType == "dm" {
		parts := strings.Split(incoming.RoomID, ":")
		if len(parts) == 3 {
			recipientID := parts[1]
			if recipientID == client.ID {
				recipientID = parts[2]
			}
			_ = h.notificationService.SendNotification(
				recipientID, "", client.ID, client.Username,
				"message", client.Username+" sent you a message",
			)
		}
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
