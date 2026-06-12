package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/utils"
)

type notifItem struct {
	ID   string `json:"id"`
	Type string `json:"type"`
	Read bool   `json:"read"`
}

func listNotifications(t *testing.T, router *gin.Engine, token string) []notifItem {
	t.Helper()
	w := authedRequest(t, router, "GET", "/api/notification", token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Data []notifItem `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp.Data
}

func findNotifByType(notifs []notifItem, notifType string) *notifItem {
	for i := range notifs {
		if notifs[i].Type == notifType {
			return &notifs[i]
		}
	}
	return nil
}

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
		Data []struct {
			Read bool `json:"read"`
		} `json:"data"`
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
	if resp.Total < 1 {
		t.Fatalf("expected the notification to stay listed after mark read, got %d", resp.Total)
	}
	for _, n := range resp.Data {
		if !n.Read {
			t.Fatalf("expected every notification read after mark all read")
		}
	}
}

func TestNotification_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "GET", "/api/notification", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestNotification_FollowNotifies(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nfol-a", "nfol-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nfol-b", "nfol-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("follow: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d", w.Code)
	}
	var resp struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total < 1 {
		t.Fatalf("expected at least 1 notification for follow, got %d", resp.Total)
	}
}

func TestNotification_UnfriendNotifies(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nunf-a", "nunf-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nunf-b", "nunf-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")

	w := authedRequest(t, router, "DELETE", "/api/friends/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("remove friend: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d", w.Code)
	}
	var resp struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total < 1 {
		t.Fatalf("expected at least 1 notification for unfriend, got %d", resp.Total)
	}
}

func TestNotification_DeleteRemovesIt(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "ndel-a", "ndel-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "ndel-b", "ndel-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")

	notifs := listNotifications(t, router, bob.Token)
	if len(notifs) == 0 {
		t.Fatalf("expected at least 1 notification before delete")
	}
	notifID := notifs[0].ID

	w := authedRequest(t, router, "DELETE", "/api/notification/"+notifID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete notification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	for _, n := range listNotifications(t, router, bob.Token) {
		if n.ID == notifID {
			t.Fatalf("expected notification %s to be deleted", notifID)
		}
	}
}

func TestNotification_DeleteMissingIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "ndelm-a", "ndelm-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/notification/"+utils.NewID(), alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete missing notification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestNotification_DeleteRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "DELETE", "/api/notification/"+utils.NewID(), "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestNotification_DeleteCannotTouchOtherUsers(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "ndelo-a", "ndelo-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "ndelo-b", "ndelo-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")

	notifs := listNotifications(t, router, bob.Token)
	if len(notifs) == 0 {
		t.Fatalf("expected at least 1 notification")
	}
	notifID := notifs[0].ID

	w := authedRequest(t, router, "DELETE", "/api/notification/"+notifID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete as other user: expected 200 no-op, got %d", w.Code)
	}

	var still bool
	for _, n := range listNotifications(t, router, bob.Token) {
		if n.ID == notifID {
			still = true
		}
	}
	if !still {
		t.Fatalf("expected bob's notification to survive a delete by alice")
	}
}

func TestNotification_FriendRequestRemovedOnAccept(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nfra-a", "nfra-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nfra-b", "nfra-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") == nil {
		t.Fatalf("expected a friend_request notification before accept")
	}

	w := authedRequest(t, router, "POST", "/api/friends/accept/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("accept: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") != nil {
		t.Fatalf("expected the friend_request notification to be removed after accept")
	}
	if findNotifByType(listNotifications(t, router, alice.Token), "friend_accept") == nil {
		t.Fatalf("expected alice to receive a friend_accept notification")
	}
}

func TestNotification_FriendRequestRemovedOnReject(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nfrr-a", "nfrr-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nfrr-b", "nfrr-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") == nil {
		t.Fatalf("expected a friend_request notification before reject")
	}

	w := authedRequest(t, router, "POST", "/api/friends/reject/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("reject: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") != nil {
		t.Fatalf("expected the friend_request notification to be removed after reject")
	}
}

func TestNotification_FriendRequestRemovedOnCancel(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nfrc-a", "nfrc-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nfrc-b", "nfrc-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") == nil {
		t.Fatalf("expected a friend_request notification before cancel")
	}

	w := authedRequest(t, router, "DELETE", "/api/friends/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("cancel request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	if findNotifByType(listNotifications(t, router, bob.Token), "friend_request") != nil {
		t.Fatalf("expected the friend_request notification to be removed after cancel")
	}
}

func TestNotification_AutoAcceptSendsFriendAccept(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "naa-a", "naa-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "naa-b", "naa-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")

	w := authedRequest(t, router, "POST", "/api/friends/request/"+alice.ID, bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("cross request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	bobNotifs := listNotifications(t, router, bob.Token)
	if findNotifByType(bobNotifs, "friend_request") != nil {
		t.Fatalf("expected bob's friend_request notification to be removed after auto-accept")
	}
	aliceNotifs := listNotifications(t, router, alice.Token)
	if findNotifByType(aliceNotifs, "friend_request") != nil {
		t.Fatalf("expected alice to get no friend_request notification on auto-accept")
	}
	if findNotifByType(aliceNotifs, "friend_accept") == nil {
		t.Fatalf("expected alice to receive a friend_accept notification on auto-accept")
	}

	w = authedRequest(t, router, "GET", "/api/friends/status/"+bob.ID, alice.Token, "")
	var status struct {
		Status string `json:"status"`
	}
	json.Unmarshal(w.Body.Bytes(), &status)
	if status.Status != "friends" {
		t.Fatalf("expected friends status after auto-accept, got %q", status.Status)
	}
}

func TestNotification_ClearAllRemovesEverything(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "nclr-a", "nclr-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "nclr-b", "nclr-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, "")
	if len(listNotifications(t, router, bob.Token)) < 2 {
		t.Fatalf("expected at least 2 notifications before clear")
	}

	w := authedRequest(t, router, "DELETE", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("clear all: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	if left := listNotifications(t, router, bob.Token); len(left) != 0 {
		t.Fatalf("expected 0 notifications after clear all, got %d", len(left))
	}
}

func TestNotification_ClearAllRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "DELETE", "/api/notification", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
