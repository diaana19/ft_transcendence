package test

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

// readUntilType reads frames until one with the given "type" arrives or a 4s
// deadline elapses.
func readUntilType(t *testing.T, conn *websocket.Conn, wantType string) map[string]any {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(4 * time.Second))
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

// wsOpen subscribes the connection to its conversation with peerID. The server
// replies with a "history" frame (asserted/drained by the caller).
func wsOpen(t *testing.T, conn *websocket.Conn, peerID string) {
	t.Helper()
	if err := conn.WriteJSON(map[string]any{"action": "open", "peer_id": peerID}); err != nil {
		t.Fatalf("write open: %v", err)
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

func TestWS_SendAndReceiveDM(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wsalice", "wsalice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "wsbob", "wsbob@test.com", "StrongPass123!")

	aliceConn := wsConnect(t, srv.URL, alice.Token)
	defer aliceConn.Close()
	bobConn := wsConnect(t, srv.URL, bob.Token)
	defer bobConn.Close()

	// Both participants open the conversation so the dm channel subscription
	// becomes live for realtime delivery.
	wsOpen(t, aliceConn, bob.ID)
	wsOpen(t, bobConn, alice.ID)
	time.Sleep(600 * time.Millisecond)

	if err := aliceConn.WriteJSON(map[string]any{
		"action":       "message",
		"recipient_id": bob.ID,
		"content":      "hi bob",
	}); err != nil {
		t.Fatalf("write message: %v", err)
	}

	msg := readUntilType(t, bobConn, "message")
	inner, ok := msg["message"].(map[string]any)
	if !ok {
		t.Fatalf("expected message payload, got %v", msg)
	}
	if inner["content"] != "hi bob" {
		t.Fatalf("expected content 'hi bob', got %v", inner["content"])
	}
	if inner["sender_id"] != alice.ID || inner["recipient_id"] != bob.ID {
		t.Fatalf("unexpected sender/recipient: %v", inner)
	}
}

func TestWS_OpenReturnsHistory(t *testing.T) {
	router, db := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	alice := registerAndLogin(t, router, "wshist_a", "wshist_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "wshist_b", "wshist_b@test.com", "StrongPass123!")

	seedMessage(t, db, alice.ID, bob.ID, "earlier message")

	conn := wsConnect(t, srv.URL, alice.Token)
	defer conn.Close()

	wsOpen(t, conn, bob.ID)

	hist := readUntilType(t, conn, "history")
	msgs, ok := hist["messages"].([]any)
	if !ok || len(msgs) != 1 {
		t.Fatalf("expected 1 history message, got %v", hist["messages"])
	}
	first, _ := msgs[0].(map[string]any)
	if first["content"] != "earlier message" {
		t.Fatalf("unexpected history content: %v", first["content"])
	}
}

func TestWS_UnknownActionAndEmptyOpen(t *testing.T) {
	router, _ := SetupTestEnv()
	srv := httptest.NewServer(router)
	defer srv.Close()

	u := registerAndLogin(t, router, "wsleave", "wsleave@test.com", "StrongPass123!")
	conn := wsConnect(t, srv.URL, u.Token)
	defer conn.Close()

	// Unknown actions and an empty/self peer open must be no-ops, not crashes.
	if err := conn.WriteJSON(map[string]any{"action": "bogus"}); err != nil {
		t.Fatalf("write bogus: %v", err)
	}
	if err := conn.WriteJSON(map[string]any{"action": "open", "peer_id": ""}); err != nil {
		t.Fatalf("write empty open: %v", err)
	}
	if err := conn.WriteJSON(map[string]any{"action": "open", "peer_id": u.ID}); err != nil {
		t.Fatalf("write self open: %v", err)
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

	msg := readUntilType(t, conn, "notification")
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

	wsOpen(t, aliceConn, bob.ID)
	wsOpen(t, bobConn, alice.ID)
	time.Sleep(500 * time.Millisecond)

	// attachment is rejected (not private) so no chat message is broadcast
	aliceConn.WriteJSON(map[string]any{"action": "message", "recipient_id": bob.ID, "file_id": publicFile})

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

	wsOpen(t, aliceConn, bob.ID)
	wsOpen(t, bobConn, alice.ID)
	time.Sleep(600 * time.Millisecond)

	if err := aliceConn.WriteJSON(map[string]any{
		"action":       "message",
		"recipient_id": bob.ID,
		"file_id":      fileID,
	}); err != nil {
		t.Fatalf("write attachment message: %v", err)
	}

	readUntilType(t, bobConn, "message")

	// after the DM attachment, bob is granted access (canAccess true → 404 from disk, not 403)
	w = authedRequest(t, router, "GET", "/api/files/"+fileID, bob.Token, "")
	if w.Code == http.StatusForbidden {
		t.Fatalf("post-share access: expected access granted, still got 403")
	}
}
