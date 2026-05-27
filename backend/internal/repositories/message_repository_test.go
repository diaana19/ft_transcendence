package repositories

import (
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

func newTestMessageDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if err := db.AutoMigrate(&models.Message{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	return db
}

// MessageRepository read/write happy paths (PollSince, ListConversation,
// Create, GetByRoomID, GetReplies) are covered end-to-end through the /chat,
// /ws and /rooms endpoints against a real Postgres. Only the query/connection
// error branches — unreachable over HTTP — are unit-tested here by closing the
// underlying connection.

func TestMessageRepo_BootstrapQueryError(t *testing.T) {
	db := newTestMessageDB(t)
	r := NewMessageRepository(db)

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("get sql.DB: %v", err)
	}
	sqlDB.Close()

	if _, err := r.PollSince("a", "", 10); err == nil {
		t.Error("expected error on bootstrap poll with closed DB")
	}
	if _, err := r.PollSince("a", "01", 10); err == nil {
		t.Error("expected error on cursor poll with closed DB")
	}
	if _, err := r.ListConversation("a", "b", "", 10); err == nil {
		t.Error("expected error on bootstrap list with closed DB")
	}
	if err := r.Create(&models.Message{ID: "x"}); err == nil {
		t.Error("expected error on create with closed DB")
	}
}
