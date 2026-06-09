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

// chatHistoryLimit is the max number of past messages loaded when a conversation is opened.
const chatHistoryLimit = 50

// ChatHandler manages the WebSocket chat. It handles the connection, direct messages,
// notifications and the Redis pub/sub used to share messages between server instances.
type ChatHandler struct {
	manager             *WSManager
	rdb                 *redis.Client
	notificationService *services.NotificationService
	msgRepo             repositories.MessageRepository
	users               repositories.UserRepository
	fileRepo            repositories.FileRepository
	upgrader            websocket.Upgrader
	subscribedRooms     map[string]bool
	subscribedMu        sync.Mutex
}

// IncomingMessage is the message the client sends over the WebSocket.
type IncomingMessage struct {
	Action      string  `json:"action"`
	PeerID      string  `json:"peer_id"`
	RecipientID string  `json:"recipient_id"`
	Content     string  `json:"content"`
	FileID      *string `json:"file_id,omitempty"`
}

// OutgoingMessage is the message the server sends back to the client over the WebSocket.
type OutgoingMessage struct {
	Type     string                   `json:"type"`
	Message  *models.MessageResponse  `json:"message,omitempty"`
	Messages []models.MessageResponse `json:"messages,omitempty"`
	PeerID   string                   `json:"peer_id,omitempty"`
}

// NewChatHandler creates a ChatHandler. It also sets the WebSocket upgrader to only
// accept connections from the allowed origins.
func NewChatHandler(
	manager *WSManager,
	rdb *redis.Client,
	notifService *services.NotificationService,
	msgRepo repositories.MessageRepository,
	users repositories.UserRepository,
	fileRepo repositories.FileRepository,
	frontendURL string,
) *ChatHandler {
	allowed := []string{frontendURL, "null", ""}
	return &ChatHandler{
		manager:             manager,
		rdb:                 rdb,
		notificationService: notifService,
		msgRepo:             msgRepo,
		users:               users,
		fileRepo:            fileRepo,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return slices.Contains(allowed, r.Header.Get("Origin"))
			},
		},
		subscribedRooms: make(map[string]bool),
	}
}

// sendPendingNotifications sends all unread notifications to the client. It marks them
// as read only when every notification was delivered.
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
	if allDelivered {
		_ = h.notificationService.MarkAllRead(client.ID)
	}
}

// @Summary      Real-time chat WebSocket
// @Description  Upgrades the connection to a WebSocket for real-time chat and presence. This is not a regular
// @Description  REST call: the client must perform a WebSocket handshake. Authentication is via the "auth_token"
// @Description  cookie or a "token" query parameter carrying the JWT access token.
// @Tags         chat
// @Param        token query string false "JWT access token (alternative to the auth_token cookie)"
// @Success      101 {string} string "Switching Protocols — WebSocket established"
// @Failure      401 {object} map[string]string "missing or invalid token"
// @Router       /ws/chat [get]
//
// HandleWS upgrades the request to a WebSocket. It registers the client, subscribes to its
// notification channel, sends pending notifications and starts the read and write pumps.
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

// readPump reads messages from the client until the connection closes. It passes each
// message to HandleMessage.
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

// HandleMessage parses the raw message and routes it by action. It supports "open" to open
// a conversation and "message" to send a direct message.
func (h *ChatHandler) HandleMessage(client *Client, raw []byte) {
	var incoming IncomingMessage
	if err := json.Unmarshal(raw, &incoming); err != nil {
		log.Printf("Invalid message format from %s: %v\n", client.ID, err)
		return
	}

	switch incoming.Action {
	case "open":
		h.handleOpen(client, incoming.PeerID)
	case "message":
		h.handleDM(client, incoming)
	default:
		log.Printf("Unknown action from %s: %s\n", client.ID, incoming.Action)
	}
}

// dmChannel builds the channel name for a direct message between two users. It always
// sorts the ids so both users get the same channel name.
func dmChannel(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return "dm:" + a + ":" + b
}

// handleOpen joins the client to the DM room, subscribes the room to Redis once, and
// sends the conversation history back to the client.
func (h *ChatHandler) handleOpen(client *Client, peerID string) {
	if peerID == "" || peerID == client.ID {
		return
	}

	channel := dmChannel(client.ID, peerID)
	h.manager.JoinRoom(client, channel)

	h.subscribedMu.Lock()
	if !h.subscribedRooms[channel] {
		h.subscribedRooms[channel] = true
		redispub.Subscribe(context.Background(), h.rdb, "chat:"+channel, func(payload string) {
			h.manager.BroadcastToRoom(channel, []byte(payload), "")
		})
	}
	h.subscribedMu.Unlock()

	msgs, err := h.msgRepo.ListConversation(client.ID, peerID, "", chatHistoryLimit)
	if err != nil {
		log.Printf("[Chat] history load failed for %s<->%s: %v", client.ID, peerID, err)
		return
	}

	out := OutgoingMessage{Type: "history", PeerID: peerID, Messages: toResponses(msgs)}
	payload, err := json.Marshal(out)
	if err != nil {
		log.Printf("Marshal error: %v\n", err)
		return
	}
	safeSend(client.Send, payload)
}

// handleDM validates and saves a direct message, then publishes it to the room and sends
// a notification to the recipient. It ignores empty or self-addressed messages.
func (h *ChatHandler) handleDM(client *Client, incoming IncomingMessage) {
	content := strings.TrimSpace(incoming.Content)
	if content == "" && incoming.FileID == nil {
		return
	}
	if incoming.RecipientID == "" || incoming.RecipientID == client.ID {
		return
	}
	if _, err := h.users.GetByID(incoming.RecipientID); err != nil {
		log.Printf("[Chat] message to unknown recipient %s: %v", incoming.RecipientID, err)
		return
	}

	if incoming.FileID != nil {
		if err := h.grantAttachment(client.ID, incoming.RecipientID, *incoming.FileID); err != nil {
			log.Printf("[Chat] attachment rejected: %v", err)
			return
		}
	}

	id, err := uuid.NewV7()
	if err != nil {
		log.Printf("uuid error: %v", err)
		return
	}

	msg := models.Message{
		ID:          id.String(),
		CreatedAt:   time.Now().UTC(),
		SenderID:    client.ID,
		Username:    client.Username,
		RecipientID: incoming.RecipientID,
		Content:     content,
		FileID:      incoming.FileID,
		Type:        "dm",
	}

	if err := h.msgRepo.Create(&msg); err != nil {
		log.Printf("Failed to save message: %v", err)
		return
	}

	resp := msg.ToResponse()
	h.publishToRoom(dmChannel(client.ID, incoming.RecipientID), OutgoingMessage{
		Type:    "message",
		Message: &resp,
	})

	_ = h.notificationService.SendNotification(
		incoming.RecipientID, "", client.ID, client.Username,
		"message", client.Username+" sent you a message",
		"",
	)
}

// toResponses converts a slice of messages to the response model sent to the client.
func toResponses(msgs []models.Message) []models.MessageResponse {
	out := make([]models.MessageResponse, len(msgs))
	for i := range msgs {
		out[i] = msgs[i].ToResponse()
	}
	return out
}

// publishToRoom marshals the message and publishes it to the room's Redis channel.
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

// grantAttachment gives the recipient access to a file sent in a DM. The file must be
// owned by the sender and uploaded as private, otherwise it returns an error.
func (h *ChatHandler) grantAttachment(senderID, recipientID, fileID string) error {
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

	if err := h.fileRepo.GrantAccess(fileID, recipientID); err != nil {
		return fmt.Errorf("failed to grant access: %w", err)
	}

	log.Printf("[Chat] attachment granted: fileID=%s, sender=%s → recipient=%s",
		fileID, senderID, recipientID)
	return nil
}
