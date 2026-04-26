package socket

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Transcendence/models"
	redispub "github.com/Transcendence/redis"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// origin := r.Header.Get("Origin")
		// return origin == "http://localhost:3000"
		return true
	},
}

type ChatHandler struct {
	manager *WSManager
	rdb     *redis.Client
}

type IncomingMessage struct {
	Action   string  `json:"action"`
	RoomID   string  `json:"room_id"`
	Content  string  `json:"content"`
	ParentID *string `json:"parent_id"`
}

type OutgoingMessage struct {
	Type    string          `json:"type"`
	Message *models.Message `json:"message,omitempty"`
	UserID  string          `json:"user_id,omitempty"`
	RoomID  string          `json:"room_id,omitempty"`
}

func NewChatHandler(manager *WSManager, rdb *redis.Client) *ChatHandler {
	return &ChatHandler{manager: manager, rdb: rdb}
}

func (h *ChatHandler) HandleWS(c *gin.Context) {

	var userID string
	if id, exists := c.Get("userID"); exists {
		userID = id.(string)
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
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v\n", err)
		return
	}

	client := &Client{
		ID:   userID,
		Conn: conn,
		Send: make(chan []byte, 256),
	}

	h.manager.RegisterClient(client)
	defer h.manager.UnregisterClient(client)

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

	redispub.Subscribe(h.rdb, "chat:"+roomID, func(payload string) {
		h.manager.BroadcastToRoom(roomID, []byte(payload), "")
	})

	out := OutgoingMessage{
		Type:   "joined",
		UserID: client.ID,
		RoomID: roomID,
	}
	h.publishToRoom(roomID, out, "")
}

func (h *ChatHandler) handleLeave(client *Client, roomID string) {
	if roomID == "" {
		return
	}
	h.manager.LeaveRoom(client, roomID)

	out := OutgoingMessage{
		Type:   "left",
		UserID: client.ID,
		RoomID: roomID,
	}
	h.publishToRoom(roomID, out, client.ID)
}

func (h *ChatHandler) handleChat(client *Client, incoming IncomingMessage) {
	if incoming.Content == "" || incoming.RoomID == "" {
		return
	}

	msg := models.Message{
		ID:        uuid.New().String(),
		CreatedAt: time.Now(),
		SenderID:  client.ID,
		RoomID:    incoming.RoomID,
		Content:   incoming.Content,
		ParentID:  incoming.ParentID,
	}

	out := OutgoingMessage{
		Type:    "message",
		Message: &msg,
	}

	h.publishToRoom(incoming.RoomID, out, client.ID)
}

func (h *ChatHandler) publishToRoom(roomID string, out OutgoingMessage, senderID string) {
	payload, err := json.Marshal(out)
	if err != nil {
		log.Printf("Marshal error: %v\n", err)
		return
	}

	if err := redispub.Publish(h.rdb, "chat:"+roomID, string(payload)); err != nil {
		log.Printf("Publish error: %v\n", err)
	}
}
