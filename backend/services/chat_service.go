package services

import (
	"errors"
	"strings"
	"time"

	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	defaultChatLimit = 50
	maxChatLimit     = 200
)

var (
	ErrEmptyContent      = errors.New("content must not be empty")
	ErrRecipientNotFound = errors.New("recipient not found")
)

var newMessageID = func() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

type ChatService struct {
	messages repositories.MessageRepository
	users    repositories.UserRepository
}

func NewChatService(messages repositories.MessageRepository, users repositories.UserRepository) *ChatService {
	return &ChatService{messages: messages, users: users}
}

func (s *ChatService) Send(senderID string, input models.CreateMessageInput) (*models.MessageResponse, error) {
	content := strings.TrimSpace(input.Content)
	if content == "" {
		return nil, ErrEmptyContent
	}

	if _, err := s.users.GetByID(input.RecipientID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrRecipientNotFound
		}
		return nil, err
	}

	id, err := newMessageID()
	if err != nil {
		return nil, err
	}

	msg := &models.Message{
		ID:          id,
		SenderID:    senderID,
		RecipientID: input.RecipientID,
		Content:     content,
		CreatedAt:   time.Now().UTC(),
	}
	if err := s.messages.Create(msg); err != nil {
		return nil, err
	}
	resp := msg.ToResponse()
	return &resp, nil
}

func (s *ChatService) Poll(userID, since string, limit int) (*models.PollResponse, error) {
	limit = clampLimit(limit)
	msgs, err := s.messages.PollSince(userID, since, limit)
	if err != nil {
		return nil, err
	}
	return buildPollResponse(msgs, since), nil
}

func (s *ChatService) ListConversation(userID, peerID, since string, limit int) (*models.PollResponse, error) {
	limit = clampLimit(limit)
	msgs, err := s.messages.ListConversation(userID, peerID, since, limit)
	if err != nil {
		return nil, err
	}
	return buildPollResponse(msgs, since), nil
}

func clampLimit(limit int) int {
	if limit <= 0 {
		return defaultChatLimit
	}
	if limit > maxChatLimit {
		return maxChatLimit
	}
	return limit
}

func buildPollResponse(msgs []models.Message, fallbackCursor string) *models.PollResponse {
	out := &models.PollResponse{
		Messages:   make([]models.MessageResponse, len(msgs)),
		NextCursor: fallbackCursor,
	}
	for i, m := range msgs {
		out.Messages[i] = m.ToResponse()
	}
	if len(msgs) > 0 {
		out.NextCursor = msgs[len(msgs)-1].ID
	}
	return out
}
