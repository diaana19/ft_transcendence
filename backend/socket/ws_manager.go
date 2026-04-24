package websocket

import (
	"log"
	"sync"

	"github.com/Transcendence/models"
)

type WSManager struct {
	mu      sync.RWMutex
	rooms   map[string]map[string]*models.Client
	clients map[string]*models.Client
}

func NewWSManager() *WSManager {
	return &WSManager{
		rooms:   make(map[string]map[string]*models.Client),
		clients: make(map[string]*models.Client),
	}
}

func (m *WSManager) RegisterClient(client *models.Client) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.clients[client.ID] = client
	log.Printf("Client %s connected\n", client.ID)
}


