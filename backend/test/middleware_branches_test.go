package test

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"ft_transcendence/backend/internal/middleware"
	"ft_transcendence/backend/internal/utils"
)

func TestRateLimitMiddleware_Branches(t *testing.T) {
	SetupTestEnv()

	// Default ceiling when RATE_LIMIT_MAX is unset/invalid (exercises the
	// fallback return in rateLimitMax).
	old, had := os.LookupEnv("RATE_LIMIT_MAX")
	os.Unsetenv("RATE_LIMIT_MAX")
	_ = middleware.RateLimitMiddleware(sharedRDB)
	if had {
		os.Setenv("RATE_LIMIT_MAX", old)
	}

	// Redis error -> 500.
	errHF := middleware.RateLimitMiddleware(brokenRDB(t))
	c, w := testCtx(http.MethodGet, "/", "")
	errHF(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("rate limit redis error: expected 500, got %d", w.Code)
	}

	// Exceeding the ceiling -> 429.
	os.Setenv("RATE_LIMIT_MAX", "1")
	limitHF := middleware.RateLimitMiddleware(sharedRDB)
	if had {
		defer os.Setenv("RATE_LIMIT_MAX", old)
	} else {
		defer os.Unsetenv("RATE_LIMIT_MAX")
	}
	c1, _ := testCtx(http.MethodGet, "/", "")
	limitHF(c1)
	c2, w2 := testCtx(http.MethodGet, "/", "")
	limitHF(c2)
	if w2.Code != http.StatusTooManyRequests {
		t.Fatalf("rate limit exceeded: expected 429, got %d", w2.Code)
	}
}

func TestOptionalAuthMiddleware_Branches(t *testing.T) {
	SetupTestEnv()
	hf := middleware.OptionalAuthMiddleware(sharedRDB)

	// Invalid token falls through to the anonymous path (no 401, no user_id).
	c, _ := testCtx(http.MethodGet, "/", "")
	c.Request.Header.Set("Authorization", "Bearer not-a-jwt")
	hf(c)
	if _, exists := c.Get("user_id"); exists {
		t.Fatal("invalid token should not set user_id")
	}

	// Blacklisted token also falls through without setting identity.
	token, err := utils.GenerateJWT(utils.NewID(), "optauth")
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}
	if err := sharedRDB.Set(context.Background(), "blacklist:"+token, "1", time.Minute).Err(); err != nil {
		t.Fatalf("blacklist set: %v", err)
	}
	c2, _ := testCtx(http.MethodGet, "/", "")
	c2.Request.Header.Set("Authorization", "Bearer "+token)
	hf(c2)
	if _, exists := c2.Get("user_id"); exists {
		t.Fatal("blacklisted token should not set user_id")
	}
}
