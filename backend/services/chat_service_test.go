package services

import (
	"errors"
	"sort"
	"strings"
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
	sort.Slice(msgs, func(i, j int) bool { return msgs[i].ID < msgs[j].ID })
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

func TestChatSend_Success(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a", Username: "alice"}
	userRepo.users["b"] = &models.User{ID: "b", Username: "bob"}

	resp, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "hello"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.SenderID != "a" || resp.RecipientID != "b" || resp.Content != "hello" {
		t.Errorf("unexpected response: %+v", resp)
	}
	if !strings.Contains(resp.ID, "-") || len(resp.ID) != 36 {
		t.Errorf("expected UUID id, got %q", resp.ID)
	}
	if len(msgRepo.messages) != 1 {
		t.Errorf("expected 1 stored message, got %d", len(msgRepo.messages))
	}
}

func TestChatSend_TrimsContent(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	userRepo.users["b"] = &models.User{ID: "b"}

	resp, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "  hi  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.Content != "hi" {
		t.Errorf("expected trimmed content 'hi', got %q", resp.Content)
	}
}

func TestChatSend_EmptyContent(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	userRepo.users["b"] = &models.User{ID: "b"}

	_, err := svc.Send("a", models.CreateMessageInput{RecipientID: "b", Content: "   "})
	if !errors.Is(err, ErrEmptyContent) {
		t.Fatalf("expected ErrEmptyContent, got %v", err)
	}
}

func TestChatSend_UnknownRecipient(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}

	_, err := svc.Send("a", models.CreateMessageInput{RecipientID: "ghost", Content: "hi"})
	if !errors.Is(err, ErrRecipientNotFound) {
		t.Fatalf("expected ErrRecipientNotFound, got %v", err)
	}
}

func TestChatPoll_BootstrapEmptyCursor(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	for _, id := range []string{"01", "02", "03"} {
		msgRepo.messages = append(msgRepo.messages, models.Message{
			ID: id, SenderID: "a", RecipientID: "x", Content: "m" + id,
		})
	}

	resp, err := svc.Poll("a", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(resp.Messages))
	}
	if resp.Messages[0].ID != "01" || resp.Messages[2].ID != "03" {
		t.Errorf("expected ascending order, got %+v", resp.Messages)
	}
	if resp.NextCursor != "03" {
		t.Errorf("expected next_cursor=03, got %q", resp.NextCursor)
	}
}

func TestChatPoll_OnlyNewerThanCursor(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	for _, id := range []string{"01", "02", "03", "04"} {
		msgRepo.messages = append(msgRepo.messages, models.Message{
			ID: id, SenderID: "a", RecipientID: "x", Content: "m" + id,
		})
	}

	resp, err := svc.Poll("a", "02", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(resp.Messages))
	}
	if resp.Messages[0].ID != "03" || resp.Messages[1].ID != "04" {
		t.Errorf("expected [03,04], got %+v", resp.Messages)
	}
	if resp.NextCursor != "04" {
		t.Errorf("expected next_cursor=04, got %q", resp.NextCursor)
	}
}

func TestChatPoll_EmptyResultKeepsCursor(t *testing.T) {
	svc, userRepo, _ := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}

	resp, err := svc.Poll("a", "07", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != 0 {
		t.Errorf("expected no messages, got %d", len(resp.Messages))
	}
	if resp.NextCursor != "07" {
		t.Errorf("expected next_cursor preserved, got %q", resp.NextCursor)
	}
}

func TestChatListConversation_ExcludesThirdParty(t *testing.T) {
	svc, userRepo, msgRepo := newChatTestService()
	userRepo.users["a"] = &models.User{ID: "a"}
	userRepo.users["b"] = &models.User{ID: "b"}
	msgRepo.messages = []models.Message{
		{ID: "01", SenderID: "a", RecipientID: "b", Content: "ab"},
		{ID: "02", SenderID: "a", RecipientID: "c", Content: "ac"},
		{ID: "03", SenderID: "b", RecipientID: "a", Content: "ba"},
	}

	resp, err := svc.ListConversation("a", "b", "", 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp.Messages) != 2 {
		t.Fatalf("expected 2 messages, got %d", len(resp.Messages))
	}
	for _, m := range resp.Messages {
		if m.Content == "ac" {
			t.Error("third-party message must be excluded")
		}
	}
}

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
