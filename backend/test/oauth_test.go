package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
)

func TestOAuthLogin_Redirects(t *testing.T) {
	// OAuthService.IsConfigured short-circuits when client_id or client_secret
	// is empty, so the test must supply non-empty values to exercise the real
	// happy path. t.Setenv auto-restores the original value when the test ends.
	t.Setenv("GITHUB_CLIENT_ID", "test-client-id")
	t.Setenv("GITHUB_CLIENT_SECRET", "test-client-secret")
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("GET", "/api/auth/oauth/github/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("oauth login: expected 307, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "github.com/login/oauth/authorize") {
		t.Fatalf("expected redirect to github authorize, got %q", loc)
	}
	if !strings.Contains(loc, "state=") {
		t.Fatalf("authorize URL must carry a state param, got %q", loc)
	}
}

// TestOAuthLogin_NotConfigured covers the fail-fast guard: when the GitHub
// credentials are not set, /api/auth/oauth/github/login redirects to the SPA
// login page with an error query param rather than handing the browser a
// github.com URL with client_id= empty (which renders as a 404 on GitHub).
func TestOAuthLogin_NotConfigured(t *testing.T) {
	t.Setenv("GITHUB_CLIENT_ID", "")
	t.Setenv("GITHUB_CLIENT_SECRET", "")
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("GET", "/api/auth/oauth/github/login", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("expected 307, got %d", w.Code)
	}
	loc := w.Header().Get("Location")
	if !strings.Contains(loc, "error=oauth_not_configured") {
		t.Fatalf("expected redirect with oauth_not_configured error, got %q", loc)
	}
}

func TestOAuthCallback_NoCode(t *testing.T) {
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("GET", "/api/auth/oauth/github/callback", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("callback no code: expected 307, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "error=oauth_denied") {
		t.Fatalf("expected oauth_denied redirect, got %q", loc)
	}
}

func TestOAuthCallback_InvalidState(t *testing.T) {
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("GET", "/api/auth/oauth/github/callback?code=abc&state=bogus-state", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("callback invalid state: expected 307, got %d", w.Code)
	}
	if loc := w.Header().Get("Location"); !strings.Contains(loc, "error=invalid_state") {
		t.Fatalf("expected invalid_state redirect, got %q", loc)
	}
}

func newOAuthService(t *testing.T) *services.OAuthService {
	t.Helper()
	cfg, _ := config.Load()
	repo := repositories.NewUserRepository(sharedDB)
	return services.NewOAuthService(repo, sharedRDB, cfg)
}

func TestOAuthService_VerifyAndConsumeState(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	// empty state is rejected without touching redis
	if ok, err := svc.VerifyAndConsumeState(ctx, ""); err != nil || ok {
		t.Fatalf("empty state: expected (false,nil), got (%v,%v)", ok, err)
	}

	// unknown state is rejected
	if ok, err := svc.VerifyAndConsumeState(ctx, "never-stored"); err != nil || ok {
		t.Fatalf("unknown state: expected (false,nil), got (%v,%v)", ok, err)
	}

	// a state produced by GenerateState verifies exactly once
	state, err := svc.GenerateState(ctx)
	if err != nil {
		t.Fatalf("generate state: %v", err)
	}
	if ok, err := svc.VerifyAndConsumeState(ctx, state); err != nil || !ok {
		t.Fatalf("valid state: expected (true,nil), got (%v,%v)", ok, err)
	}
	if ok, _ := svc.VerifyAndConsumeState(ctx, state); ok {
		t.Fatal("state must be single-use")
	}
}

func TestOAuthService_BuildAuthURL(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	url := svc.BuildAuthURL("my-state")
	if !strings.Contains(url, "state=my-state") {
		t.Fatalf("auth url missing state: %s", url)
	}
}

func TestOAuthService_ExchangeCodeForToken_EmptyCode(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	if _, err := svc.ExchangeCodeForToken(context.Background(), ""); err == nil {
		t.Fatal("expected error for empty authorization code")
	}
}

func TestOAuthService_FindOrCreateUser(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	gh := &services.GitHubUser{
		ID:        time.Now().UnixNano() % 1000000000,
		Login:     "ghnewuser",
		Name:      "GH New User",
		Email:     "ghnew@example.com",
		AvatarURL: "https://avatar.example/x.png",
	}

	created, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("find or create (new): %v", err)
	}
	if created.Provider != "github" || created.Username == "" {
		t.Fatalf("unexpected created user: %+v", created)
	}

	// second call with the same github id returns the existing user
	again, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("find or create (existing by github id): %v", err)
	}
	if again.ID != created.ID {
		t.Fatalf("expected same user id, got %s vs %s", again.ID, created.ID)
	}
}

func TestOAuthService_FindOrCreateUser_LinksExistingEmail(t *testing.T) {
	router, _ := SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	// a local account already exists with this email
	registerAndLogin(t, router, "linklocal", "linkme@example.com", "StrongPass123!")

	gh := &services.GitHubUser{
		ID:    time.Now().UnixNano()%1000000000 + 1,
		Login: "linklocal",
		Name:  "Link Local",
		Email: "linkme@example.com",
	}

	linked, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("link existing email: %v", err)
	}
	if linked.Email != "linkme@example.com" {
		t.Fatalf("expected linked account to keep email, got %s", linked.Email)
	}
	if linked.GithubID == nil {
		t.Fatal("expected github id to be linked onto existing account")
	}
}

func TestOAuthService_FindOrCreateUser_UsernameCollision(t *testing.T) {
	router, _ := SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	// occupy the desired username with a local user
	registerAndLogin(t, router, "octocat", "octocat@example.com", "StrongPass123!")

	gh := &services.GitHubUser{
		ID:    time.Now().UnixNano()%1000000000 + 2,
		Login: "octocat",
		Email: "freshoctocat@example.com",
	}

	user, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("username collision: %v", err)
	}
	if user.Username == "octocat" {
		t.Fatalf("expected a de-duplicated username, got %q", user.Username)
	}
	_ = *user
}
