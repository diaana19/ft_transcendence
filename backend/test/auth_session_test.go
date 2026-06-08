package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestAuth_RefreshToken(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "refreshu", "refreshu@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/refresh", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("refresh: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Fatal("expected a token in refresh response")
	}
}

func TestAuth_RefreshMissingToken(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/refresh", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestAuth_LogoutRevokesToken(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "logoutu", "logoutu@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/logout", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("logout: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/users", u.Token, "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 after logout, got %d", w.Code)
	}
}
