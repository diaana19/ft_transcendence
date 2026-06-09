package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/oauth2"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/redis"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestOAuth_FindOrCreateUser_NewAndExisting(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	gh := &services.GitHubUser{ID: 4242, Login: "octocat", Name: "The Octocat", Email: "octo@github.test", AvatarURL: "http://avatar"}

	created, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("create from github: %v", err)
	}
	if created.Username != "octocat" || created.Provider != "github" {
		t.Fatalf("unexpected created user: %+v", created)
	}

	again, err := svc.FindOrCreateUser(ctx, gh)
	if err != nil {
		t.Fatalf("lookup by github id: %v", err)
	}
	if again.ID != created.ID {
		t.Fatalf("expected same user on second call, got %s vs %s", again.ID, created.ID)
	}
}

func TestOAuth_FindOrCreateUser_UsernameCollision(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	first, err := svc.FindOrCreateUser(ctx, &services.GitHubUser{ID: 1, Login: "dup", Email: "dup1@github.test"})
	if err != nil {
		t.Fatalf("first user: %v", err)
	}
	second, err := svc.FindOrCreateUser(ctx, &services.GitHubUser{ID: 2, Login: "dup", Email: "dup2@github.test"})
	if err != nil {
		t.Fatalf("second user: %v", err)
	}
	if first.Username != "dup" || second.Username != "dup1" {
		t.Fatalf("expected collision resolution dup/dup1, got %q/%q", first.Username, second.Username)
	}
}

func TestOAuth_FindOrCreateUser_LinkByEmail(t *testing.T) {
	router, db := SetupTestEnv()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	registerAndLogin(t, router, "linkme", "linkme@github.test", "StrongPass123!")

	svc := services.NewOAuthService(repositories.NewUserRepository(db), sharedRDB, cfg)

	linked, err := svc.FindOrCreateUser(context.Background(), &services.GitHubUser{ID: 99, Login: "linkme", Email: "linkme@github.test"})
	if err != nil {
		t.Fatalf("link by email: %v", err)
	}
	if linked.GithubID == nil || *linked.GithubID != "99" {
		t.Fatalf("expected github id linked, got %+v", linked.GithubID)
	}
}

func TestOAuth_StateLifecycle(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	ctx := context.Background()

	state, err := svc.GenerateState(ctx)
	if err != nil {
		t.Fatalf("GenerateState: %v", err)
	}
	if state == "" {
		t.Fatal("expected non-empty state")
	}

	ok, err := svc.VerifyAndConsumeState(ctx, state)
	if err != nil || !ok {
		t.Fatalf("expected state valid, got ok=%v err=%v", ok, err)
	}

	ok, _ = svc.VerifyAndConsumeState(ctx, state)
	if ok {
		t.Fatal("state should be single-use")
	}

	ok, _ = svc.VerifyAndConsumeState(ctx, "")
	if ok {
		t.Fatal("empty state must be invalid")
	}
}

func TestOAuth_BuildAuthURLAndConfigured(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)
	if url := svc.BuildAuthURL("xyz"); url == "" {
		t.Fatal("expected non-empty auth url")
	}
	_ = svc.IsConfigured()

	if _, err := svc.ExchangeCodeForToken(context.Background(), ""); err == nil {
		t.Fatal("empty code should error")
	}
}

func TestOAuth_ExchangeAndCallbackFailure(t *testing.T) {
	router, _ := SetupTestEnv()
	svc := newOAuthService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	if _, err := svc.ExchangeCodeForToken(ctx, "bogus-code"); err == nil {
		t.Fatal("expected exchange to fail for bogus authorization code")
	}

	state, err := svc.GenerateState(context.Background())
	if err != nil {
		t.Fatalf("GenerateState: %v", err)
	}

	req, _ := http.NewRequestWithContext(ctx, "GET", "/api/auth/oauth/github/callback?code=bogus&state="+state, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("callback with valid state but bad code: expected 307, got %d - %s", w.Code, w.Body.String())
	}
}

func TestServices_NotFoundBranches(t *testing.T) {
	SetupTestEnv()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	repo := repositories.NewUserRepository(sharedDB)

	twofa := services.NewTwoFAService(repo)
	if _, gerr := twofa.GenerateSecret(utils.NewID()); gerr == nil {
		t.Fatal("GenerateSecret for unknown user should error")
	}

	usvc := services.NewUserService(repo)
	if verr := usvc.VerifyPassword(utils.NewID(), "whatever"); verr == nil {
		t.Fatal("VerifyPassword for unknown user should error")
	}

	oauth := services.NewOAuthService(repo, sharedRDB, cfg)
	ghUser, err := oauth.FindOrCreateUser(context.Background(), &services.GitHubUser{ID: 7777, Login: "nopw", Email: "nopw@gh.test"})
	if err != nil {
		t.Fatalf("FindOrCreateUser: %v", err)
	}
	if err := usvc.VerifyPassword(ghUser.ID, "anything"); err == nil {
		t.Fatal("VerifyPassword for password-less oauth account should error")
	}
}

func TestOAuth_FetchGitHubUser_BadToken(t *testing.T) {
	SetupTestEnv()
	svc := newOAuthService(t)

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	tok := &oauth2.Token{AccessToken: "ghp_thisisdefinitelynotavalidtoken00000000"}
	if _, err := svc.FetchGitHubUser(ctx, tok); err == nil {
		t.Fatal("expected error fetching github user with an invalid token")
	}
}

func TestRedis_RoundTrip(t *testing.T) {
	SetupTestEnv()
	if err := redis.RoundTrip(context.Background(), sharedRDB, "test:roundtrip", "ping-payload"); err != nil {
		t.Fatalf("RoundTrip: %v", err)
	}
	if err := redis.Publish(sharedRDB, "test:roundtrip", "lone-publish"); err != nil {
		t.Fatalf("Publish: %v", err)
	}
}
