package socket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

type WSManager struct {
	mu      sync.RWMutex
	rooms   map[string]map[string]*Client
	clients map[string]*Client
}

func NewWSManager() *WSManager {
	return &WSManager{
		rooms:   make(map[string]map[string]*Client),
		clients: make(map[string]*Client),
	}
}

func (m *WSManager) RegisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.ID] = client
	log.Printf("Client %s connected\n", client.Username)
}

func (m *WSManager) UnregisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for roomID, room := range m.rooms {
		if _, ok := room[client.ID]; ok {
			delete(room, client.ID)
			log.Printf("Client %s has left the room [%s]\n", client.ID, roomID)
			if len(room) == 0 {
				delete(m.rooms, roomID)
			}
		}
	}
	delete(m.clients, client.ID)
	close(client.Send)
	log.Printf("Client %s disconnect\n", client.ID)
}

func (m *WSManager) JoinRoom(client *Client, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[string]*Client)
	}
	m.rooms[roomID][client.ID] = client
	log.Printf("Client %s has joined room [%s]\n", client.Username, roomID)
}

func (m *WSManager) LeaveRoom(client *Client, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if room, ok := m.rooms[roomID]; ok {
		delete(room, client.ID)
		log.Printf("Client %s has left the room [%s]\n", client.Username, roomID)
		if len(room) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

// safeSend writes msg to ch without blocking. Returns true on enqueue,
// false if the buffer is full or the channel has been closed. Recover guards
// the send-on-closed-channel panic so callers can keep iterating other clients.
func safeSend(ch chan []byte, msg []byte) (delivered bool) {
	defer func() {
		if r := recover(); r != nil {
			delivered = false
			log.Printf("safeSend: channel was closed, dropping message\n")
		}
	}()
	select {
	case ch <- msg:
		return true
	default:
		log.Printf("safeSend: buffer full, dropping message\n")
		return false
	}
}

// BroadcastToRoom fans message out to every client in roomID except senderID.
// The room snapshot is taken under RLock and the actual sends happen after the
// lock is released, so a slow consumer can never block Register/Unregister/
// Join/Leave on the manager.
func (m *WSManager) BroadcastToRoom(roomID string, message []byte, senderID string) {
	m.mu.RLock()
	room, ok := m.rooms[roomID]
	if !ok {
		m.mu.RUnlock()
		return
	}
	clients := make([]*Client, 0, len(room))
	for _, c := range room {
		if c.ID == senderID {
			continue
		}
		clients = append(clients, c)
	}
	m.mu.RUnlock()

	for _, c := range clients {
		safeSend(c.Send, message)
	}
}

func (m *WSManager) GetRoomMembers(roomID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	members := make([]string, 0, len(m.rooms[roomID]))
	for clientID := range m.rooms[roomID] {
		members = append(members, clientID)
	}
	return members
}

// IsOnline returns true if the user has an active WebSocket connection.
func (m *WSManager) IsOnline(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.clients[userID]
	return ok
}

func (c *Client) WritePump() {
	defer func() { _ = c.Conn.Close() }()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error for client %s: %v\n", c.ID, err)
			return
		}
	}
}
