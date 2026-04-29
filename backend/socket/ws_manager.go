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
		if _, ok := room[roomID]; ok {
			delete(room, client.ID)
			log.Printf("Client %s has left the room [%s\n]", client.ID, roomID)
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
		log.Printf("Client %s has left the room [%s\n]", client.Username, roomID)
		if len(room) == 0 {
			delete(m.rooms, roomID)
		}
	}
}

func safeSend(ch chan []byte, msg []byte) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("safeSend: channel was closed, dropping message\n")
		}
	}()
	select {
	case ch <- msg:
	default:
		log.Printf("safeSend: buffer full, dropping message\n")
	}
}

func (m *WSManager) BroadcastToRoom(roomID string, message []byte, senderID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	room, ok := m.rooms[roomID]
	if !ok {
		return
	}
	for _, client := range room {
		if client.ID == senderID {
			continue
		}
		safeSend(client.Send, message)
	}
}

func (m *WSManager) GetRoomMembers(roomID string) []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	members := []string{}
	for clientID := range m.rooms[roomID] {
		members = append(members, clientID)
	}
	return members
}

func (c *Client) WritePump() {
	defer c.Conn.Close()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Printf("Write error for client %s: %v\n", c.ID, err)
			return
		}
	}
}
