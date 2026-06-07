package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestOnlineStatus_NotOnline(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "online1", "online1@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+user.ID+"/online", "", "")
	if w.Code != http.StatusOK {
		t.Fatalf("online status: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	var resp struct {
		Online bool `json:"online"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Online {
		t.Fatal("expected offline for user without WebSocket connection")
	}
}

func TestOnlineStatus_UnknownUser(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/users/550e8400-e29b-41d4-a716-446655440000/online", "", "")
	if w.Code != http.StatusOK {
		t.Fatalf("online status unknown user: expected 200, got %d", w.Code)
	}

	var resp struct {
		Online bool `json:"online"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Online {
		t.Fatal("expected offline for unknown user")
	}
}
