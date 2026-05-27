package test

import (
	"net/http"
	"testing"
)

const fakeUUID = "550e8400-e29b-41d4-a716-446655440000"

func TestMsg_GetRoomMessagesEmpty(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "msgu", "msgu@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/rooms/"+fakeUUID+"/messages", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("room messages: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestMsg_GetRepliesEmpty(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "msgr", "msgr@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/messages/"+fakeUUID+"/replies", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("replies: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestMsg_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "GET", "/api/rooms/"+fakeUUID+"/messages", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
