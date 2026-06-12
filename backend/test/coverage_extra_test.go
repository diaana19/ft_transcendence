package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNotification_MarkSingleRead(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "msralice", "msralice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "msrbob", "msrbob@test.com", "StrongPass123!")

	if w := authedRequest(t, router, "POST", "/api/friends/request/"+bob.ID, alice.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("friend request: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w := authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d", w.Code)
	}
	var list struct {
		Data []struct {
			ID   string `json:"id"`
			Read bool   `json:"read"`
		} `json:"data"`
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &list)
	if list.Total < 1 || len(list.Data) == 0 {
		t.Fatalf("expected at least 1 notification, got %d", list.Total)
	}
	notifID := list.Data[0].ID

	w = authedRequest(t, router, "PATCH", "/api/notification/"+notifID+"/read", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("mark single read: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	json.Unmarshal(w.Body.Bytes(), &list)
	if list.Total < 1 {
		t.Fatalf("expected the notification to stay listed after mark read, got %d", list.Total)
	}
	for _, n := range list.Data {
		if n.ID == notifID && !n.Read {
			t.Fatalf("expected notification %s to be read", notifID)
		}
	}
}

func TestNotification_MarkSingleReadMissingIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "msrnf", "msrnf@test.com", "StrongPass123!")

	fakeID := "550e8400-e29b-41d4-a716-446655440000"
	w := authedRequest(t, router, "PATCH", "/api/notification/"+fakeID+"/read", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("mark unknown notification read: expected 200 no-op, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestNotification_MarkSingleReadRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "PATCH", "/api/notification/whatever/read", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestGetUsers_ByUsername_Success(t *testing.T) {
	router, _ := SetupTestEnv()
	target := registerAndLogin(t, router, "lookupme", "lookupme@test.com", "StrongPass123!")
	caller := registerAndLogin(t, router, "lookupcaller", "lookupcaller@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users?username=lookupme", caller.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("username lookup: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	var user struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &user); err != nil {
		t.Fatalf("decode single-user response: %v", err)
	}
	if user.ID != target.ID || user.Username != "lookupme" {
		t.Fatalf("expected the looked-up user (id=%s), got id=%s username=%s", target.ID, user.ID, user.Username)
	}
}

func TestGetUsers_ByUsername_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	caller := registerAndLogin(t, router, "lookupnf", "lookupnf@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users?username=ghost-no-such-user", caller.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("unknown username lookup: expected 404, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestRefreshToken_Success(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "refresher", "refresher@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/refresh", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Fatal("refresh: expected a non-empty token")
	}
}

func TestRefreshToken_MissingToken(t *testing.T) {
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("POST", "/api/auth/refresh", &bytes.Buffer{})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("refresh without token: expected 401, got %d", w.Code)
	}
}

func TestRefreshToken_InvalidToken(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/refresh", "not-a-real-jwt", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("refresh with invalid token: expected 401, got %d", w.Code)
	}
}

func TestLogout_ViaCookie(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cookielogout", "cookielogout@test.com", "StrongPass123!")

	req, _ := http.NewRequest("POST", "/api/auth/logout", &bytes.Buffer{})
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: user.Token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("cookie logout: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestLogout_InvalidToken(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/logout", "garbage-token", "")
	if w.Code != http.StatusOK {
		t.Fatalf("logout with invalid token: expected 200, got %d", w.Code)
	}
}
