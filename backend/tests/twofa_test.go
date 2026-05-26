package tests

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/pquerna/otp/totp"
)

func TestTwoFA_SetupEnableDisableFlow(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfa", "tfa@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("setup: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var setup struct {
		Secret string `json:"secret"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &setup); err != nil || setup.Secret == "" {
		t.Fatalf("setup: bad secret response: %v body=%s", err, w.Body.String())
	}

	code, err := totp.GenerateCode(setup.Secret, time.Now())
	if err != nil {
		t.Fatalf("generate code: %v", err)
	}
	w = authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"`+code+`"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("enable: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	code2, _ := totp.GenerateCode(setup.Secret, time.Now())
	w = authedRequest(t, router, "POST", "/api/2fa/disable", u.Token, `{"code":"`+code2+`"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("disable: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_EnableInvalidCode(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfabad", "tfabad@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	w := authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"000000"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("enable invalid: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_SetupRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "POST", "/api/2fa/setup", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
