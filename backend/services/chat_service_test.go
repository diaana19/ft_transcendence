package services

import (
	"errors"
	"slices"
	"testing"

	"github.com/Transcendence/models"
)

type mockMessageRepository struct {
	messages []models.Message
	err      error
}

func newMockMessageRepo() *mockMessageRepository {
	return &mockMessageRepository{}
}

func (m *mockMessageRepository) Create(msg *models.Message) error {
	if m.err != nil {
		return m.err
	}
	m.messages = append(m.messages, *msg)
	return nil
}

func (m *mockMessageRepository) PollSince(userID, since string, limit int) ([]models.Message, error) {
	if m.err != nil {
		return nil, m.err
	}
	var matched []models.Message
	for _, msg := range m.messages {
		if msg.SenderID != userID && msg.RecipientID != userID {
			continue
		}
		matched = append(matched, msg)
	}
	return cursorSlice(matched, since, limit), nil
}

func (m *mockMessageRepository) GetByRoomID(_, _ string, _ int) ([]models.Message, error) {
	return nil, nil
}

func (m *mockMessageRepository) GetReplies(_ string) ([]models.Message, error) {
	return nil, nil
}

func (m *mockMessageRepository) ListConversation(userID, peerID, since string, limit int) ([]models.Message, error) {
	if m.err != nil {
		return nil, m.err
	}
	var matched []models.Message
	for _, msg := range m.messages {
		ab := msg.SenderID == userID && msg.RecipientID == peerID
		ba := msg.SenderID == peerID && msg.RecipientID == userID
		if ab || ba {
			matched = append(matched, msg)
		}
	}
	return cursorSlice(matched, since, limit), nil
}

func cursorSlice(msgs []models.Message, since string, limit int) []models.Message {
	slices.SortFunc(msgs, func(a, b models.Message) int {
		switch {
		case a.ID < b.ID:
			return -1
		case a.ID > b.ID:
			return 1
		default:
			return 0
		}
	})
	if since == "" {
		if len(msgs) > limit {
			msgs = msgs[len(msgs)-limit:]
		}
		return msgs
	}
	var out []models.Message
	for _, m := range msgs {
		if m.ID > since {
			out = append(out, m)
			if len(out) == limit {
				break
			}
		}
	}
	return out
}

func newChatTestService() (*ChatService, *mockUserRepository, *mockMessageRepository) {
	userRepo := newMockRepo()
	msgRepo := newMockMessageRepo()
	svc := NewChatService(msgRepo, userRepo)
	return svc, userRepo, msgRepo
}

// ChatService happy paths (send, content trimming, empty-content rejection,
// unknown recipient, poll bootstrap/cursor, conversation filtering) are covered
// end-to-end through the /chat endpoints. The tests below target branches that
// the HTTP layer cannot reach: forced repository/DB errors, a forced UUID
// generation failure, and limit clamping (which would need hundreds of seeded
// messages to observe over the API).

func TestChatPoll_LimitClamping(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	for i := range 250 {
		msgRepo.messages = append(msgRepo.messages, models.Message{
			ID:          padID(i),
			SenderID:    "a",
			RecipientID: "x",
		})
	}

	resp, err := svc.Poll("a", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != defaultChatLimit {
		t.Errorf("default limit: expected %d, got %d", defaultChatLimit, len(resp.Messages))
	}

	resp, err = svc.Poll("a", "", 5000)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != maxChatLimit {
		t.Errorf("max limit: expected %d, got %d", maxChatLimit, len(resp.Messages))
	}
}

func TestChatSend_UserRepoGenericError(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.err = errors.New("db down")

	_, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "hi"})
	if err == nil || errors.Is(err, ErrRecipientNotFound) {
		t.Fatalf("expected raw db error, got %v", err)
	}
}

func TestChatSend_UUIDError(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.users["b"] = &models.User{ID: "b"}

	original := newMessageID
	defer func() { newMessageID = original }()
	newMessageID = func() (string, error) { return "", errors.New("rng failure") }

	_, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "hi"})
	if err == nil || err.Error() != "rng failure" {
		t.Fatalf("expected rng failure, got %v", err)
	}
}

func TestChatSend_RepoCreateError(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["b"] = &models.User{ID: "b"}
	msgRepo.err = errors.New("insert fail")

	_, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "hi"})
	if err == nil || err.Error() != "insert fail" {
		t.Fatalf("expected insert fail, got %v", err)
	}
}

func TestChatPoll_RepoError(t *testing.T) {
	svc, _, msgRepo := newChatTestService()
	msgRepo.err = errors.New("query fail")

	_, err := svc.Poll("a", "", 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChatListConversation_RepoError(t *testing.T) {
	svc, _, msgRepo := newChatTestService()
	msgRepo.err = errors.New("query fail")

	_, err := svc.ListConversation("a", "b", "", 0)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestChatPoll_LimitWithinBounds(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	for i := range 10 {
		msgRepo.messages = append(msgRepo.messages, models.Message{
			ID: padID(i + 1), SenderID: "a", RecipientID: "x",
		})
	}

	resp, err := svc.Poll("a", "", 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != 5 {
		t.Errorf("expected 5 (within-bounds limit), got %d", len(resp.Messages))
	}
}

func padID(i int) string {
	const width = 6
	s := ""
	for n := i; n > 0 || s == ""; n /= 10 {
		s = string(rune('0'+n%10)) + s
		if n == 0 {
			break
		}
	}
	for len(s) < width {
		s = "0" + s
	}
	return s
}
