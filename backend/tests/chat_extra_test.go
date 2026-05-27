package tests

import (
	"fmt"
	"net/http"
	"testing"
)

func TestChat_SendEmptyContentRejected(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "cesend_a", "cesend_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "cesend_b", "cesend_b@test.com", "StrongPass123!")

	body := fmt.Sprintf(`{"recipient_id":"%s","content":"   "}`, bob.ID)
	w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty content: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestChat_SendTrimsContent(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "cetrim_a", "cetrim_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "cetrim_b", "cetrim_b@test.com", "StrongPass123!")

	body := fmt.Sprintf(`{"recipient_id":"%s","content":"  padded  "}`, bob.ID)
	w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, body)
	if w.Code != http.StatusCreated {
		t.Fatalf("send: expected 201, got %d", w.Code)
	}
	if body := w.Body.String(); !contains(body, `"content":"padded"`) {
		t.Fatalf("expected trimmed content, got %s", body)
	}
}

func TestChat_PollWithExplicitLimit(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "cepoll", "cepoll@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/chat/poll?limit=5", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("poll with limit: expected 200, got %d", w.Code)
	}

	// a non-numeric limit falls back to the default rather than erroring
	w = authedRequest(t, router, "GET", "/api/chat/poll?limit=notanumber", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("poll with invalid limit: expected 200, got %d", w.Code)
	}
}

func TestChat_SendInvalidJSON(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "cebind", "cebind@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/chat/messages", u.Token, `{"recipient_id":}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid json: expected 400, got %d", w.Code)
	}
}
