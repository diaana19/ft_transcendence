package test

import (
	"testing"
	"time"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

func TestMessageRepository_GetByRoomID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewMessageRepository(db)

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "msguser1", Email: "msg1@test.com"}
	user2 := models.User{ID: user2ID, Username: "msguser2", Email: "msg2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Create messages
	roomID := "dm:" + user1ID + ":" + user2ID
	msg1 := models.Message{ID: utils.NewID(), RoomID: roomID, SenderID: user1ID, Content: "hello", CreatedAt: time.Now()}
	msg2 := models.Message{ID: utils.NewID(), RoomID: roomID, SenderID: user2ID, Content: "hi", CreatedAt: time.Now().Add(time.Second)}
	repo.Create(&msg1)
	repo.Create(&msg2)

	// Get by room ID
	msgs, err := repo.GetByRoomID(roomID, "", 10)
	if err != nil {
		t.Fatalf("GetByRoomID: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
}

func TestMessageRepository_GetByRoomIDEmpty(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewMessageRepository(db)

	msgs, err := repo.GetByRoomID("nonexistent", "", 10)
	if err != nil {
		t.Fatalf("GetByRoomID empty: %v", err)
	}
	if len(msgs) != 0 {
		t.Fatalf("expected 0 messages, got %d", len(msgs))
	}
}

func TestMessageRepository_GetReplies(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewMessageRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "replyuser", Email: "reply@test.com"}
	db.Create(&user)

	// Create parent message
	parentID := utils.NewID()
	parent := models.Message{ID: parentID, RoomID: "room1", SenderID: userID, Content: "parent", CreatedAt: time.Now()}
	repo.Create(&parent)

	// Create reply
	reply := models.Message{ID: utils.NewID(), RoomID: "room1", SenderID: userID, Content: "reply", ParentID: &parentID, CreatedAt: time.Now().Add(time.Second)}
	repo.Create(&reply)

	// Get replies
	replies, err := repo.GetReplies(parentID)
	if err != nil {
		t.Fatalf("GetReplies: %v", err)
	}
	if len(replies) != 1 {
		t.Fatalf("expected 1 reply, got %d", len(replies))
	}
}

func TestMessageRepository_PollSince(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewMessageRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "polluser", Email: "poll@test.com"}
	db.Create(&user)

	// Create message
	msg := models.Message{ID: utils.NewID(), RoomID: "room1", SenderID: userID, Content: "test", CreatedAt: time.Now()}
	repo.Create(&msg)

	// Poll since before the message
	msgs, err := repo.PollSince(userID, "", 10)
	if err != nil {
		t.Fatalf("PollSince: %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("expected 1 message, got %d", len(msgs))
	}
}

func TestMessageRepository_ListConversation(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewMessageRepository(db)

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "convuser1", Email: "conv1@test.com"}
	user2 := models.User{ID: user2ID, Username: "convuser2", Email: "conv2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Create messages
	msg1 := models.Message{ID: utils.NewID(), SenderID: user1ID, RecipientID: user2ID, Content: "hello", CreatedAt: time.Now()}
	msg2 := models.Message{ID: utils.NewID(), SenderID: user2ID, RecipientID: user1ID, Content: "hi", CreatedAt: time.Now().Add(time.Second)}
	repo.Create(&msg1)
	repo.Create(&msg2)

	// List conversation
	msgs, err := repo.ListConversation(user1ID, user2ID, "", 10)
	if err != nil {
		t.Fatalf("ListConversation: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
}
