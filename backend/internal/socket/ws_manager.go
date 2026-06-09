package socket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client is one connected WebSocket user. Send is the buffered channel for outgoing messages.
type Client struct {
	ID       string
	Username string
	Conn     *websocket.Conn
	Send     chan []byte
}

// WSManager keeps track of connected clients and the rooms they are in. It is safe for
// concurrent use.
type WSManager struct {
	mu      sync.RWMutex
	rooms   map[string]map[string]*Client
	clients map[string]*Client
}

// NewWSManager creates an empty WSManager.
func NewWSManager() *WSManager {
	return &WSManager{
		rooms:   make(map[string]map[string]*Client),
		clients: make(map[string]*Client),
	}
}

// RegisterClient adds the client to the manager so it is marked as online.
func (m *WSManager) RegisterClient(client *Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.ID] = client
	log.Printf("Client %s connected\n", client.Username)
}

// UnregisterClient removes the client from all rooms and from the manager, and closes its
// Send channel.
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

// JoinRoom adds the client to the room, creating the room when it does not exist.
func (m *WSManager) JoinRoom(client *Client, roomID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.rooms[roomID] == nil {
		m.rooms[roomID] = make(map[string]*Client)
	}
	m.rooms[roomID][client.ID] = client
	log.Printf("Client %s has joined room [%s]\n", client.Username, roomID)
}

// LeaveRoom removes the client from the room and deletes the room when it becomes empty.
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

// safeSend sends a message to the channel without blocking. It returns false and drops the
// message when the channel is closed or its buffer is full.
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

// BroadcastToRoom sends the message to every client in the room, except the sender.
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

// GetRoomMembers returns the ids of all clients in the room.
func (m *WSManager) GetRoomMembers(roomID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	members := make([]string, 0, len(m.rooms[roomID]))
	for clientID := range m.rooms[roomID] {
		members = append(members, clientID)
	}
	return members
}

// IsOnline returns true when the user has an active connection.
func (m *WSManager) IsOnline(userID string) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()

	_, ok := m.clients[userID]
	return ok
}

// WritePump reads from the client's Send channel and writes each message to the connection.
// It returns and closes the connection when a write fails or the channel is closed.
func (c *Client) WritePump() {
	defer func() { _ = c.Conn.Close() }()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error for client %s: %v\n", c.ID, err)
			return
		}
	}
}
