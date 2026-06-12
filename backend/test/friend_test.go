package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestFriend_FollowAndListFollowers(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "falice", "falice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "fbob", "fbob@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("follow: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/users/"+bob.ID+"/followers", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("followers: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 follower, got %d", len(resp.Data))
	}

	w = authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/following", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("following: expected 200, got %d", w.Code)
	}
}

func TestFriend_FollowSelfIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "fself", "fself@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+alice.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("follow self: expected 200, got %d", w.Code)
	}
}

func TestFriend_Unfollow(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "funf-a", "funf-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "funf-b", "funf-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")

	w := authedRequest(t, router, "DELETE", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("unfollow: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestFriend_RequestAcceptFlow(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "freq-a", "freq-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "freq-b", "freq-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("send request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("accept request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/friends", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("friends list: expected 200, got %d", w.Code)
	}
}

func TestFriend_RequestUnknownTarget(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "funk", "funk@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/request/550e8400-e29b-41d4-a716-446655440000", alice.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("request to unknown: expected 400, got %d", w.Code)
	}
}

func TestFriend_StatusLifecycle(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "fst-a", "fst-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "fst-b", "fst-b@test.com", "StrongPass123!")

	getStatus := func(token, targetID string) string {
		w := authedRequest(t, router, "GET", "/api/friends/status/"+targetID, token, "")
		if w.Code != http.StatusOK {
			t.Fatalf("status: expected 200, got %d - body: %s", w.Code, w.Body.String())
		}
		var resp struct {
			Status string `json:"status"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		return resp.Status
	}

	if s := getStatus(alice.Token, bob.ID); s != "none" {
		t.Fatalf("expected none, got %q", s)
	}

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if s := getStatus(alice.Token, bob.ID); s != "pending_sent" {
		t.Fatalf("expected pending_sent, got %q", s)
	}
	if s := getStatus(bob.Token, alice.ID); s != "pending_received" {
		t.Fatalf("expected pending_received, got %q", s)
	}

	authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")
	if s := getStatus(alice.Token, bob.ID); s != "friends" {
		t.Fatalf("expected friends, got %q", s)
	}

	authedRequest(t, router, "DELETE", "/api/friends/"+bob.ID, alice.Token, "")
	if s := getStatus(alice.Token, bob.ID); s != "none" {
		t.Fatalf("expected none after remove, got %q", s)
	}
}

func TestFriend_RejectRequest(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "frej-a", "frej-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "frej-b", "frej-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")

	w := authedRequest(t, router, "POST", "/api/friends/reject/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("reject: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}
