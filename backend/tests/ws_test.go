package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// wsConnect dials the chat WebSocket endpoint on srv authenticating with token.
func wsConnect(t *testing.T, srvURL, token string) *websocket.Conn {
	t.Helper()
	wsURL := "ws" + strings.TrimPrefix(srvURL, "http") + "/api/ws/chat?token=" + token
	conn, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		code := 0
		if resp != nil {
			code = resp.StatusCode
		}
		t.Fatalf("ws dial failed: %v (status %d)", err, code)
	}
	return conn
}

// readUntilType reads frames until one with the given "type" arrives or timeout.
func readUntilType(t *testing.T, conn *websocket.Conn, wantType string, timeout time.Duration) map[string]any {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(timeout))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			t.Fatalf("read ws (waiting for %q): %v", wantType, err)
		}
		var msg map[string]any
		if err := json.Unmarshal(raw, &msg); err != nil {
			continue
		}
		if msg["type"] == wantType {
			return msg
		}
	}
}

func wsJoin(t *testing.T, conn *websocket.Conn, room string) {
	t.Helper()
	if err := conn.WriteJSON(map[string]any{"action": "join", "room_id": room}); err != nil {
		t.Fatalf("write join: %v", err)
	}
}

func TestWS_RequiresToken(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws/chat"
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Fatal("expected handshake to fail without token")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		got := 0
		if resp != nil {
			got = resp.StatusCode
		}
		t.Fatalf("expected 401, got %d", got)
	}
}

func TestWS_InvalidToken(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	wsURL := "ws" + strings.TrimPrefix(srv.URL, "http") + "/api/ws/chat?token=not-a-valid-jwt"
	_, resp, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err == nil {
		t.Fatal("expected handshake to fail with invalid token")
	}
	if resp == nil || resp.StatusCode != http.StatusUnauthorized {
		got := 0
		if resp != nil {
			got = resp.StatusCode
		}
		t.Fatalf("expected 401, got %d", got)
	}
}

func TestWS_JoinAndBroadcastMessage(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wsalice", "wsalice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "wsbob", "wsbob@test.com", "StrongPass123!")

	aliceConn := wsConnect(t, srv.URL, alice.Token)
	defer aliceConn.Close()
	bobConn := wsConnect(t, srv.URL, bob.Token)
	defer bobConn.Close()

	room := "dm:" + alice.ID + ":" + bob.ID
	wsJoin(t, aliceConn, room)
	wsJoin(t, bobConn, room)

	// Let both joins register and the redis subscription become live.
	time.Sleep(600 * time.Millisecond)

	if err := aliceConn.WriteJSON(map[string]any{
		"action":  "message",
		"room_id": room,
		"content": "hi bob",
	}); err != nil {
		t.Fatalf("write message: %v", err)
	}

	msg := readUntilType(t, bobConn, "message", 4*time.Second)
	inner, ok := msg["message"].(map[string]any)
	if !ok {
		t.Fatalf("expected message payload, got %v", msg)
	}
	if inner["content"] != "hi bob" {
		t.Fatalf("expected content 'hi bob', got %v", inner["content"])
	}

	// The message must also be retrievable through the REST history endpoint.
	w := authedRequest(t, router, "GET", "/api/rooms/"+room+"/messages", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("room history: expected 200, got %d", w.Code)
	}
}

func TestWS_LeaveRoomAndUnknownAction(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	u := registerAndLogin(t, router, "wsleave", "wsleave@test.com", "StrongPass123!")
	conn := wsConnect(t, srv.URL, u.Token)
	defer conn.Close()

	room := "dm:" + u.ID + ":someone"
	wsJoin(t, conn, room)
	time.Sleep(300 * time.Millisecond)

	if err := conn.WriteJSON(map[string]any{"action": "leave", "room_id": room}); err != nil {
		t.Fatalf("write leave: %v", err)
	}
	if err := conn.WriteJSON(map[string]any{"action": "bogus"}); err != nil {
		t.Fatalf("write bogus: %v", err)
	}
	// empty room id join is a no-op and must not crash
	if err := conn.WriteJSON(map[string]any{"action": "join", "room_id": ""}); err != nil {
		t.Fatalf("write empty join: %v", err)
	}
	time.Sleep(200 * time.Millisecond)
}

func TestWS_DeliversPendingNotificationsOnConnect(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	sender := registerAndLogin(t, router, "wsnotif_s", "wsnotif_s@test.com", "StrongPass123!")
	target := registerAndLogin(t, router, "wsnotif_t", "wsnotif_t@test.com", "StrongPass123!")

	// a friend request creates an unread notification for the target
	w := authedRequest(t, router, "POST", "/api/friends/request/"+target.ID, sender.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("friend request: expected 200, got %d", w.Code)
	}

	// when the target connects, the pending notification is pushed down the socket
	conn := wsConnect(t, srv.URL, target.Token)
	defer conn.Close()

	msg := readUntilType(t, conn, "notification", 4*time.Second)
	if _, ok := msg["notification"]; !ok {
		t.Fatalf("expected a notification payload, got %v", msg)
	}
}

func TestWS_AttachmentRejectedWhenNotPrivate(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wsrej_a", "wsrej_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "wsrej_b", "wsrej_b@test.com", "StrongPass123!")

	// public file cannot be used as a DM attachment (must be private)
	publicFile := uploadAndGetID(t, router, alice.Token, "public")

	aliceConn := wsConnect(t, srv.URL, alice.Token)
	defer aliceConn.Close()
	bobConn := wsConnect(t, srv.URL, bob.Token)
	defer bobConn.Close()

	room := "dm:" + alice.ID + ":" + bob.ID
	wsJoin(t, aliceConn, room)
	wsJoin(t, bobConn, room)
	time.Sleep(500 * time.Millisecond)

	// attachment is rejected (not private) so no chat message is broadcast
	aliceConn.WriteJSON(map[string]any{"action": "message", "room_id": room, "file_id": publicFile})

	bobConn.SetReadDeadline(time.Now().Add(1 * time.Second))
	for {
		_, raw, err := bobConn.ReadMessage()
		if err != nil {
			break
		}
		var msg map[string]any
		json.Unmarshal(raw, &msg)
		if msg["type"] == "message" {
			t.Fatal("expected no broadcast for a rejected (non-private) attachment")
		}
	}
}

func TestWS_AttachmentRejectedForBadRoom(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wsbad_a", "wsbad_a@test.com", "StrongPass123!")

	privateFile := uploadAndGetID(t, router, alice.Token, "private")

	conn := wsConnect(t, srv.URL, alice.Token)
	defer conn.Close()

	// a non-"dm:a:b" room id is not a valid attachment target
	room := "group-room-123"
	wsJoin(t, conn, room)
	time.Sleep(400 * time.Millisecond)

	conn.WriteJSON(map[string]any{"action": "message", "room_id": room, "file_id": privateFile})

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			break
		}
		var msg map[string]any
		json.Unmarshal(raw, &msg)
		if msg["type"] == "message" {
			t.Fatal("attachment in a non-dm room must be rejected (no message broadcast)")
		}
	}
}

func TestWS_AttachmentGrantsAccessToRecipient(t *testing.T) {
	t.Cleanup(func() { os.RemoveAll("./uploads") })
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wsatt_a", "wsatt_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "wsatt_b", "wsatt_b@test.com", "StrongPass123!")

	fileID := uploadAndGetID(t, router, alice.Token, "private")

	// before sharing, bob cannot access alice's private file
	w := authedRequest(t, router, "GET", "/api/files/"+fileID, bob.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("pre-share access: expected 403, got %d", w.Code)
	}

	aliceConn := wsConnect(t, srv.URL, alice.Token)
	defer aliceConn.Close()
	bobConn := wsConnect(t, srv.URL, bob.Token)
	defer bobConn.Close()

	room := "dm:" + alice.ID + ":" + bob.ID
	wsJoin(t, aliceConn, room)
	wsJoin(t, bobConn, room)
	time.Sleep(600 * time.Millisecond)

	if err := aliceConn.WriteJSON(map[string]any{
		"action":  "message",
		"room_id": room,
		"file_id": fileID,
	}); err != nil {
		t.Fatalf("write attachment message: %v", err)
	}

	readUntilType(t, bobConn, "message", 4*time.Second)

	// after the DM attachment, bob is granted access (canAccess true → 404 from disk, not 403)
	w = authedRequest(t, router, "GET", "/api/files/"+fileID, bob.Token, "")
	if w.Code == http.StatusForbidden {
		t.Fatalf("post-share access: expected access granted, still got 403")
	}
}
