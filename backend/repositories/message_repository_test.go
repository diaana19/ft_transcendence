package repositories

import (
	"testing"

	"github.com/Transcendence/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

func TestMessageRepo_CreateAndPollSince(t *testing.T) {
	db := newTestMessageDB(t)
	r := NewMessageRepository(db)

	for _, id := range []string{"01", "02", "03"} {
		err := r.Create(&models.Message{ID: id, SenderID: "a", RecipientID: "b", Content: id})
		if err != nil {
			t.Fatalf("create: %v", err)
		}
	}

	msgs, err := r.PollSince("a", "01", 10)
	if err != nil {
		t.Fatalf("poll: %v", err)
	}
	if len(msgs) != 2 || msgs[0].ID != "02" || msgs[1].ID != "03" {
		t.Errorf("expected [02,03], got %+v", msgs)
	}

	msgs, err = r.PollSince("a", "", 10)
	if err != nil {
		t.Fatalf("bootstrap poll: %v", err)
	}
	if len(msgs) != 3 || msgs[0].ID != "01" || msgs[2].ID != "03" {
		t.Errorf("bootstrap expected ascending [01,02,03], got %+v", msgs)
	}
}

func TestMessageRepo_ListConversation(t *testing.T) {
	db := newTestMessageDB(t)
	r := NewMessageRepository(db)

	r.Create(&models.Message{ID: "01", SenderID: "a", RecipientID: "b", Content: "ab"})
	r.Create(&models.Message{ID: "02", SenderID: "a", RecipientID: "c", Content: "ac"})
	r.Create(&models.Message{ID: "03", SenderID: "b", RecipientID: "a", Content: "ba"})

	msgs, err := r.ListConversation("a", "b", "", 10)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(msgs) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(msgs))
	}
	for _, m := range msgs {
		if m.Content == "ac" {
			t.Error("third-party included")
		}
	}

	msgs, err = r.ListConversation("a", "b", "01", 10)
	if err != nil {
		t.Fatalf("list since: %v", err)
	}
	if len(msgs) != 1 || msgs[0].ID != "03" {
		t.Errorf("expected [03], got %+v", msgs)
	}
}

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
