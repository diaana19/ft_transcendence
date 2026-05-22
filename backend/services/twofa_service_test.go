package services

import (
	"testing"
	"time"

	"github.com/Transcendence/models"
	"github.com/pquerna/otp/totp"
)

// makeSecret returns a fresh TOTP secret usable by the 2FA service tests.
func makeSecret(t *testing.T) string {
	t.Helper()
	key, err := totp.Generate(totp.GenerateOpts{Issuer: "Transcendence", AccountName: "u@test.com"})
	if err != nil {
		t.Fatalf("generate secret: %v", err)
	}
	return key.Secret()
}

func TestTwoFAService_GenerateSecret(t *testing.T) {
	repo := newMockRepo()
	s := NewTwoFAService(repo)

	if _, err := s.GenerateSecret("missing"); err == nil {
		t.Fatal("missing user should error")
	}

	repo.users["enabled"] = &models.User{ID: "enabled", Email: "e@test.com", TwoFAEnabled: true}
	if _, err := s.GenerateSecret("enabled"); err == nil {
		t.Fatal("already enabled should error")
	}

	repo.users["u1"] = &models.User{ID: "u1", Email: "u1@test.com"}
	resp, err := s.GenerateSecret("u1")
	if err != nil || resp.Secret == "" || resp.QRCodeURL == "" {
		t.Fatalf("generate: %v %+v", err, resp)
	}
}

func TestTwoFAService_EnableTwoFA(t *testing.T) {
	repo := newMockRepo()
	s := NewTwoFAService(repo)

	if err := s.EnableTwoFA("missing", "123456"); err == nil {
		t.Fatal("missing user should error")
	}

	repo.users["nosecret"] = &models.User{ID: "nosecret"}
	if err := s.EnableTwoFA("nosecret", "123456"); err == nil {
		t.Fatal("no secret in progress should error")
	}

	secret := makeSecret(t)
	repo.users["enabled"] = &models.User{ID: "enabled", TwoFASecret: &secret, TwoFAEnabled: true}
	if err := s.EnableTwoFA("enabled", "123456"); err == nil {
		t.Fatal("already enabled should error")
	}

	repo.users["u1"] = &models.User{ID: "u1", TwoFASecret: &secret}
	if err := s.EnableTwoFA("u1", "000000"); err == nil {
		t.Fatal("invalid code should error")
	}
	code, _ := totp.GenerateCode(secret, time.Now())
	if err := s.EnableTwoFA("u1", code); err != nil {
		t.Fatalf("valid code enable: %v", err)
	}
}

func TestTwoFAService_ValidateCode(t *testing.T) {
	repo := newMockRepo()
	s := NewTwoFAService(repo)

	if _, err := s.ValidateCode("missing", "123456"); err == nil {
		t.Fatal("missing user should error")
	}

	repo.users["disabled"] = &models.User{ID: "disabled"}
	if _, err := s.ValidateCode("disabled", "123456"); err == nil {
		t.Fatal("2FA not enabled should error")
	}

	secret := makeSecret(t)
	repo.users["u1"] = &models.User{ID: "u1", TwoFASecret: &secret, TwoFAEnabled: true}
	code, _ := totp.GenerateCode(secret, time.Now())
	ok, err := s.ValidateCode("u1", code)
	if err != nil || !ok {
		t.Fatalf("valid code: ok=%v err=%v", ok, err)
	}
}

func TestTwoFAService_DisableTwoFA(t *testing.T) {
	repo := newMockRepo()
	s := NewTwoFAService(repo)
	repo.users["u1"] = &models.User{ID: "u1"}
	if err := s.DisableTwoFA("u1"); err != nil {
		t.Fatalf("disable: %v", err)
	}
}
