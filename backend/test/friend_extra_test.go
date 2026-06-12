package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestFriend_RemoveFriend(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "rmf-a", "rmf-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "rmf-b", "rmf-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")

	w := authedRequest(t, router, "DELETE", "/api/friends/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("remove friend: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/friends", alice.Token, "")
	var resp struct {
		Data []any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 0 {
		t.Fatalf("expected no friends after removal, got %d", len(resp.Data))
	}
}

func TestFriend_RemoveNonFriendIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "rmnf-a", "rmnf-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "rmnf-b", "rmnf-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/friends/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("remove non-friend: expected 200, got %d", w.Code)
	}
}

func TestFriend_DuplicateFollowIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "dupf-a", "dupf-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "dupf-b", "dupf-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("first follow: expected 200, got %d", w.Code)
	}
	w = authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("duplicate follow: expected 200, got %d", w.Code)
	}
}

func TestFriend_UnfollowNotFollowingIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "unf2-a", "unf2-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "unf2-b", "unf2-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("unfollow when not following: expected 200, got %d", w.Code)
	}
}

func TestFriend_AcceptWithoutPendingFails(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "acc-a", "acc-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "acc-b", "acc-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("accept without pending: expected 400, got %d", w.Code)
	}
}

func TestFriend_DuplicateRequestIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "dreq-a", "dreq-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "dreq-b", "dreq-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("first request: expected 200, got %d", w.Code)
	}
	w = authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("duplicate request: expected 200, got %d", w.Code)
	}
}

func TestFriend_RejectWithoutPendingIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "rej2-a", "rej2-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "rej2-b", "rej2-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/reject/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("reject without pending: expected 200, got %d", w.Code)
	}
}

func TestFriend_ListingEmpty(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "lst-a", "lst-a@test.com", "StrongPass123!")

	for _, path := range []string{"/followers", "/following", "/friends"} {
		w := authedRequest(t, router, "GET", "/api/users/"+alice.ID+path, alice.Token, "")
		if w.Code != http.StatusOK {
			t.Fatalf("listing %s: expected 200, got %d", path, w.Code)
		}
		var resp struct {
			Data []any `json:"data"`
		}
		json.Unmarshal(w.Body.Bytes(), &resp)
		if len(resp.Data) != 0 {
			t.Fatalf("listing %s: expected empty, got %d", path, len(resp.Data))
		}
	}
}
