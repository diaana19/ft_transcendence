package test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestLogin_ByUsername(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "loginuname", "loginuname@test.com", "StrongPass123!")

	body := `{"username":"loginuname","password":"StrongPass123!"}`
	w := authedRequest(t, router, "POST", "/api/auth/login", "", body)
	if w.Code != http.StatusOK {
		t.Fatalf("login by username: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Token == "" {
		t.Fatal("expected token on username login")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "loginbaduser", "loginwrong@test.com", "StrongPass123!")

	body := `{"email":"loginwrong@test.com","password":"WrongPass123!"}`
	w := authedRequest(t, router, "POST", "/api/auth/login", "", body)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("wrong password: expected 401, got %d", w.Code)
	}
}

func TestLogin_MissingIdentifier(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"password":"StrongPass123!"}`
	w := authedRequest(t, router, "POST", "/api/auth/login", "", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing identifier: expected 400, got %d", w.Code)
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/login", "", `{"email":`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid json: expected 400, got %d", w.Code)
	}
}

func TestLogin_UnknownUser(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"email":"nobody@test.com","password":"StrongPass123!"}`
	w := authedRequest(t, router, "POST", "/api/auth/login", "", body)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unknown user: expected 401, got %d", w.Code)
	}
}

func TestVerify2FA_FullFlow(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "v2fa", "v2fa@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("setup: expected 200, got %d", w.Code)
	}
	var setup struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w.Body.Bytes(), &setup)

	code, _ := totp.GenerateCode(setup.Secret, time.Now())
	w = authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, fmt.Sprintf(`{"code":"%s"}`, code))
	if w.Code != http.StatusOK {
		t.Fatalf("enable: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	loginBody := `{"email":"v2fa@test.com","password":"StrongPass123!"}`
	w = authedRequest(t, router, "POST", "/api/auth/login", "", loginBody)
	if w.Code != http.StatusOK {
		t.Fatalf("login with 2fa: expected 200, got %d", w.Code)
	}
	var loginResp struct {
		Needs2FA     bool   `json:"needs_2fa"`
		PendingToken string `json:"pending_token"`
	}
	json.Unmarshal(w.Body.Bytes(), &loginResp)
	if !loginResp.Needs2FA || loginResp.PendingToken == "" {
		t.Fatalf("expected needs_2fa with pending token, got %+v", loginResp)
	}

	verifyCode, _ := totp.GenerateCode(setup.Secret, time.Now())
	verifyBody := fmt.Sprintf(`{"pending_token":"%s","code":"%s"}`, loginResp.PendingToken, verifyCode)
	w = authedRequest(t, router, "POST", "/api/auth/2fa/verify", "", verifyBody)
	if w.Code != http.StatusOK {
		t.Fatalf("verify: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var verifyResp struct {
		Token string `json:"token"`
	}
	json.Unmarshal(w.Body.Bytes(), &verifyResp)
	if verifyResp.Token == "" {
		t.Fatal("expected a session token after 2FA verification")
	}
}

func TestVerify2FA_InvalidPendingToken(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"pending_token":"not-a-real-token","code":"123456"}`
	w := authedRequest(t, router, "POST", "/api/auth/2fa/verify", "", body)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("invalid pending token: expected 401, got %d", w.Code)
	}
}

func TestVerify2FA_InvalidBody(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/2fa/verify", "", `{"code":"123"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid body: expected 400, got %d", w.Code)
	}
}

func TestVerify2FA_WrongCode(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "v2fawc", "v2fawc@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	var setup struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w.Body.Bytes(), &setup)
	code, _ := totp.GenerateCode(setup.Secret, time.Now())
	authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, fmt.Sprintf(`{"code":"%s"}`, code))

	loginBody := `{"email":"v2fawc@test.com","password":"StrongPass123!"}`
	w = authedRequest(t, router, "POST", "/api/auth/login", "", loginBody)
	var loginResp struct {
		PendingToken string `json:"pending_token"`
	}
	json.Unmarshal(w.Body.Bytes(), &loginResp)

	verifyBody := fmt.Sprintf(`{"pending_token":"%s","code":"000000"}`, loginResp.PendingToken)
	w = authedRequest(t, router, "POST", "/api/auth/2fa/verify", "", verifyBody)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("wrong 2fa code: expected 401, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestLogout_MissingToken(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/logout", "bad.token.here", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("logout with invalid token: expected 401, got %d", w.Code)
	}
}

func TestAuthMe(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "authme", "authme@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/auth/me", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.User.ID != u.ID {
		t.Fatalf("expected user.id=%s, got %s", u.ID, resp.User.ID)
	}
	if resp.User.Username != "authme" {
		t.Fatalf("expected user.username=authme, got %q", resp.User.Username)
	}
}

func TestAuthMe_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "GET", "/api/auth/me", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d - body: %s", w.Code, w.Body.String())
	}
}
