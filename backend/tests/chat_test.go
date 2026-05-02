package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

type chatTestUser struct {
	ID    string
	Token string
}

func registerAndLogin(t *testing.T, router *gin.Engine, username, email, password string) chatTestUser {
	t.Helper()

	regBody := fmt.Sprintf(`{
		"username": "%s",
		"email": "%s",
		"password": "%s",
		"dateOfBirth": "2000-01-01"
	}`, username, email, password)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register %s: expected 200, got %d - body: %s", username, w.Code, w.Body.String())
	}
	var regResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &regResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login %s: expected 200, got %d - body: %s", username, w.Code, w.Body.String())
	}
	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	return chatTestUser{ID: loginResp.User.ID, Token: loginResp.Token}
}

func authedRequest(t *testing.T, router *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body == "" {
		buf = &bytes.Buffer{}
	} else {
		buf = bytes.NewBufferString(body)
	}
	req, _ := http.NewRequest(method, path, buf)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestChat_SendThenPollRoundTrip(t *testing.T) {
	router, _ := SetupTestEnv()

	alice := registerAndLogin(t, router, "alice", "alice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "bob", "bob@test.com", "StrongPass123!")

	sendBody := fmt.Sprintf(`{"recipient_id":"%s","content":"hello bob"}`, bob.ID)
	w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, sendBody)
	if w.Code != http.StatusCreated {
		t.Fatalf("send: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
	var sent struct {
		ID          string `json:"id"`
		SenderID    string `json:"sender_id"`
		RecipientID string `json:"recipient_id"`
		Content     string `json:"content"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &sent); err != nil {
		t.Fatalf("decode send response: %v", err)
	}
	if sent.SenderID != alice.ID || sent.RecipientID != bob.ID || sent.Content != "hello bob" {
		t.Fatalf("unexpected sent message: %+v", sent)
	}
	if len(sent.ID) != 36 {
		t.Fatalf("expected UUID id (36 chars), got %q", sent.ID)
	}

	w = authedRequest(t, router, "GET", "/api/chat/poll", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("poll: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var poll struct {
		Messages []struct {
			ID      string `json:"id"`
			Content string `json:"content"`
		} `json:"messages"`
		NextCursor string `json:"next_cursor"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &poll); err != nil {
		t.Fatalf("decode poll response: %v", err)
	}
	if len(poll.Messages) != 1 || poll.Messages[0].ID != sent.ID {
		t.Fatalf("expected one message with id %s, got %+v", sent.ID, poll.Messages)
	}
	if poll.NextCursor != sent.ID {
		t.Fatalf("expected next_cursor=%s, got %s", sent.ID, poll.NextCursor)
	}

	w = authedRequest(t, router, "GET", "/api/chat/poll?since="+poll.NextCursor, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("second poll: expected 200, got %d", w.Code)
	}
	var poll2 struct {
		Messages   []any  `json:"messages"`
		NextCursor string `json:"next_cursor"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &poll2); err != nil {
		t.Fatalf("decode second poll: %v", err)
	}
	if len(poll2.Messages) != 0 {
		t.Fatalf("expected no new messages on second poll, got %d", len(poll2.Messages))
	}
	if poll2.NextCursor != sent.ID {
		t.Fatalf("expected cursor preserved, got %s", poll2.NextCursor)
	}
}

func TestChat_UnknownRecipient(t *testing.T) {
	router, _ := SetupTestEnv()

	alice := registerAndLogin(t, router, "alice2", "alice2@test.com", "StrongPass123!")

	body := `{"recipient_id":"00000000-0000-0000-0000-000000000000","content":"hi"}`
	w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, body)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 for unknown recipient, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestChat_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"recipient_id":"x","content":"hi"}`
	req, _ := http.NewRequest("POST", "/api/chat/messages", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestChat_ListConversation_RequiresWith(t *testing.T) {
	router, _ := SetupTestEnv()

	alice := registerAndLogin(t, router, "alice3", "alice3@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/chat/messages", alice.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 without 'with' param, got %d", w.Code)
	}
}

func TestChat_ListConversation_OnlyBetweenPair(t *testing.T) {
	router, _ := SetupTestEnv()

	alice := registerAndLogin(t, router, "alice4", "alice4@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "bob4", "bob4@test.com", "StrongPass123!")
	eve := registerAndLogin(t, router, "eve4", "eve4@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/chat/messages", alice.Token,
		fmt.Sprintf(`{"recipient_id":"%s","content":"to bob"}`, bob.ID))
	authedRequest(t, router, "POST", "/api/chat/messages", alice.Token,
		fmt.Sprintf(`{"recipient_id":"%s","content":"to eve"}`, eve.ID))
	authedRequest(t, router, "POST", "/api/chat/messages", bob.Token,
		fmt.Sprintf(`{"recipient_id":"%s","content":"reply"}`, alice.ID))

	w := authedRequest(t, router, "GET", "/api/chat/messages?with="+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Messages []struct {
			Content string `json:"content"`
		} `json:"messages"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Messages) != 2 {
		t.Fatalf("expected 2 messages between alice and bob, got %d", len(resp.Messages))
	}
	for _, m := range resp.Messages {
		if m.Content == "to eve" {
			t.Fatal("conversation must not include third-party messages")
		}
	}
}
