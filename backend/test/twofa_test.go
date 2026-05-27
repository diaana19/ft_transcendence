package test

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

func TestTwoFA_SetupUnknownUser(t *testing.T) {
	router, _ := SetupTestEnv()
	// a validly-signed token whose user id has no row → service GetByID fails
	w := authedRequest(t, router, "POST", "/api/2fa/setup", tokenFor("550e8400-e29b-41d4-a716-446655440000"), "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("setup unknown user: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_SetupAlreadyEnabled(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfaalready", "tfaalready@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	var setup struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w.Body.Bytes(), &setup)
	code, _ := totp.GenerateCode(setup.Secret, time.Now())
	authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"`+code+`"}`)

	// a second setup while already enabled must be rejected
	w = authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("setup when already enabled: expected 400, got %d", w.Code)
	}
}

func TestTwoFA_EnableWithoutSetup(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfanosetup", "tfanosetup@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"123456"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("enable without setup: expected 400, got %d", w.Code)
	}
}

func TestTwoFA_EnableAlreadyEnabled(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfareenable", "tfareenable@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	var setup struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w.Body.Bytes(), &setup)
	code, _ := totp.GenerateCode(setup.Secret, time.Now())
	authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"`+code+`"}`)

	code2, _ := totp.GenerateCode(setup.Secret, time.Now())
	w = authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"`+code2+`"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("re-enable: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_EnableRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "POST", "/api/2fa/enable", "", `{"code":"123456"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("enable without auth: expected 401, got %d", w.Code)
	}
}

func TestTwoFA_EnableBadCodeFormat(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfabadfmt", "tfabadfmt@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"abc"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad code format: expected 400, got %d", w.Code)
	}
}

func TestTwoFA_DisableNotEnabled(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfadisnone", "tfadisnone@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/disable", u.Token, `{"code":"123456"}`)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("disable when not enabled: expected 500, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_DisableWrongCode(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfadiscode", "tfadiswrong@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/setup", u.Token, "")
	var setup struct {
		Secret string `json:"secret"`
	}
	json.Unmarshal(w.Body.Bytes(), &setup)
	code, _ := totp.GenerateCode(setup.Secret, time.Now())
	authedRequest(t, router, "POST", "/api/2fa/enable", u.Token, `{"code":"`+code+`"}`)

	w = authedRequest(t, router, "POST", "/api/2fa/disable", u.Token, `{"code":"000000"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("disable wrong code: expected 400, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestTwoFA_DisableRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "POST", "/api/2fa/disable", "", `{"code":"123456"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("disable without auth: expected 401, got %d", w.Code)
	}
}

func TestTwoFA_EnableUnknownUser(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "POST", "/api/2fa/enable",
		tokenFor("550e8400-e29b-41d4-a716-446655440000"), `{"code":"123456"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("enable unknown user: expected 400, got %d", w.Code)
	}
}

func TestTwoFA_DisableUnknownUser(t *testing.T) {
	router, _ := SetupTestEnv()
	w := authedRequest(t, router, "POST", "/api/2fa/disable",
		tokenFor("550e8400-e29b-41d4-a716-446655440000"), `{"code":"123456"}`)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("disable unknown user: expected 500, got %d", w.Code)
	}
}

func TestTwoFA_DisableBadCodeFormat(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "tfadisfmt", "tfadisfmt@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/2fa/disable", u.Token, `{"code":"xx"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("disable bad code format: expected 400, got %d", w.Code)
	}
}
