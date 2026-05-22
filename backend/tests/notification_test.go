package tests

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestNotification_FriendRequestNotifies(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nalice", "nalice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nbob", "nbob@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("friend request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total < 1 {
		t.Fatalf("expected at least 1 notification, got %d", resp.Total)
	}

	w = authedRequest(t, router, "PATCH", "/api/notification/read", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("mark read: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 0 {
		t.Fatalf("expected 0 unread after mark read, got %d", resp.Total)
	}
}

func TestNotification_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "GET", "/api/notification", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
