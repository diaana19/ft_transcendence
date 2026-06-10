package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// MessageRepository handles the chat messages in the database.
type MessageRepository interface {
	// Create saves a new message.
	Create(message *models.Message) error
	// GetByRoomID returns the top level messages of a room after the since cursor.
	GetByRoomID(roomID string, since string, limit int) ([]models.Message, error)
	// GetReplies returns all the replies of a message.
	GetReplies(parentID string) ([]models.Message, error)
	// PollSince returns the messages of the user after the since cursor.
	PollSince(userID, since string, limit int) ([]models.Message, error)
	// ListConversation returns the messages between two users after the since cursor.
	ListConversation(userID, peerID, since string, limit int) ([]models.Message, error)
	// ListConversationPartnerIDs returns the distinct user ids the user has direct messages with.
	ListConversationPartnerIDs(userID string) ([]string, error)
}

type messageRepository struct {
	db *gorm.DB
}

// NewMessageRepository creates a new MessageRepository using the given database.
func NewMessageRepository(db *gorm.DB) MessageRepository {
	return &messageRepository{db: db}
}

// Create saves a new message.
func (r *messageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

// GetByRoomID returns the top level messages of a room after the since cursor.
func (r *messageRepository) GetByRoomID(roomID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where("room_id = ? AND parent_id IS NULL", roomID)
	return runCursorQuery(q, since, limit)
}

// GetReplies returns all the replies of a message.
func (r *messageRepository) GetReplies(parentID string) ([]models.Message, error) {
	var replies []models.Message
	err := r.db.Where("parent_id = ?", parentID).
		Order("created_at asc").
		Find(&replies).Error
	return replies, err
}

// PollSince returns the messages of the user after the since cursor.
func (r *messageRepository) PollSince(userID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where("sender_id = ? OR recipient_id = ?", userID, userID)
	return runCursorQuery(q, since, limit)
}

// ListConversation returns the messages between two users after the since cursor.
func (r *messageRepository) ListConversation(userID, peerID, since string, limit int) ([]models.Message, error) {
	q := r.db.Where(
		"(sender_id = ? AND recipient_id = ?) OR (sender_id = ? AND recipient_id = ?)",
		userID, peerID, peerID, userID,
	)
	return runCursorQuery(q, since, limit)
}

// ListConversationPartnerIDs returns the distinct user ids the user has direct messages with.
func (r *messageRepository) ListConversationPartnerIDs(userID string) ([]string, error) {
	var ids []string
	err := r.db.Raw(`
		SELECT DISTINCT CASE WHEN sender_id = ? THEN recipient_id ELSE sender_id END AS peer_id
		FROM messages
		WHERE (sender_id = ? OR recipient_id = ?)
		  AND recipient_id IS NOT NULL AND recipient_id <> ''`,
		userID, userID, userID).Scan(&ids).Error
	return ids, err
}

// runCursorQuery runs the query with cursor pagination. When since is empty it
// returns the last messages in ascending order, else the ones after the cursor.
func runCursorQuery(q *gorm.DB, since string, limit int) ([]models.Message, error) {
	var messages []models.Message
	if since == "" {
		if err := q.Order("id DESC").Limit(limit).Find(&messages).Error; err != nil {
			return nil, err
		}
		for i, j := 0, len(messages)-1; i < j; i, j = i+1, j-1 {
			messages[i], messages[j] = messages[j], messages[i]
		}
		return messages, nil
	}
	err := q.Where("id > ?", since).Order("id ASC").Limit(limit).Find(&messages).Error
	return messages, err
}
