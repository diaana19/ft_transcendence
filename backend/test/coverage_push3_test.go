package test

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/routes"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// brokenRDB returns a redis client whose connection is closed, so every command
// fails — for exercising the redis-error branches of services.
func brokenRDB(t *testing.T) *goredis.Client {
	t.Helper()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	rdb, err := cfg.Redis.Connect()
	if err != nil {
		t.Fatalf("redis connect: %v", err)
	}
	_ = rdb.Close()
	return rdb
}

// makeFileHeader builds an in-memory *multipart.FileHeader from raw bytes.
func makeFileHeader(t *testing.T, filename string, content []byte) *multipart.FileHeader {
	t.Helper()
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	fw, err := mw.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	_, _ = fw.Write(content)
	mw.Close()

	r := multipart.NewReader(&body, mw.Boundary())
	form, err := r.ReadForm(1 << 20)
	if err != nil {
		t.Fatalf("read form: %v", err)
	}
	return form.File["file"][0]
}

// ---------------- Auth service: error branches ----------------

func TestAuthService_ErrorBranches(t *testing.T) {
	SetupTestEnv()

	// Password missing -> error (use the shared DB so the unique checks pass).
	_, db := SetupTestEnv()
	svc := services.NewAuthService(repositories.NewUserRepository(db))
	noPass := &models.User{Username: "nopass", Email: "nopass@test.com"}
	if _, err := svc.CreateAuthUserService(noPass); err == nil {
		t.Fatal("expected error for missing password")
	}

	// CreateUser fails against a broken DB.
	bsvc := services.NewAuthService(repositories.NewUserRepository(brokenDB(t)))
	pw := "StrongPass123!"
	u := &models.User{Username: "brk", Email: "brk@test.com", Password: &pw}
	if _, err := bsvc.CreateAuthUserService(u); err == nil {
		t.Fatal("expected error creating user on broken db")
	}

	// Login for an account with no local password (OAuth-only) -> invalid.
	oauthUser := &models.User{ID: utils.NewID(), Username: "ghonly", Email: "ghonly@test.com", Provider: "github"}
	db.Create(oauthUser)
	if _, err := svc.LoginAuthUserService("ghonly", "whatever"); err == nil {
		t.Fatal("expected invalid credential for passwordless account")
	}

	// Redis-backed helpers against a closed client.
	brdb := brokenRDB(t)
	if err := svc.LogoutAuthUserService("tok", time.Minute, brdb); err == nil {
		t.Fatal("expected logout error on closed redis")
	}
	if _, err := svc.CreatePendingLogin("u1", brdb); err == nil {
		t.Fatal("expected pending login error on closed redis")
	}
	if _, err := svc.PeekPendingLogin("tok", brdb); err == nil {
		t.Fatal("expected peek error on closed redis")
	}
}

// ---------------- Notification service: error branches ----------------

func TestNotificationService_ErrorBranches(t *testing.T) {
	_, db := SetupTestEnv()

	// Broken repo: GetUsernameByID warns, Create returns the error.
	brokenRepo := repositories.NewNotificationRepositories(brokenDB(t))
	bsvc := services.NewNotificationService(brokenRepo, repositories.NewNotificationPubSub(sharedRDB))
	if err := bsvc.SendNotification("u1", "", "actor", "actorname", "follow", "hi", ""); err == nil {
		t.Fatal("expected error saving notification on broken db")
	}
	// SendCommentNotification also fetches username (warns) then fails on save,
	// and exercises the >50 char preview truncation.
	long := "this is a very long comment preview that exceeds the fifty character cap easily"
	if err := bsvc.SendCommentNotification("u1", "actor", "actorname", "p1", long); err == nil {
		t.Fatal("expected error from SendCommentNotification on broken db")
	}

	// Real repo but a broken pub/sub client: Create succeeds, publish logs error.
	user := &models.User{ID: utils.NewID(), Username: "notifuser", Email: "notif@test.com"}
	db.Create(user)
	goodRepo := repositories.NewNotificationRepositories(db)
	pubBroken := services.NewNotificationService(goodRepo, repositories.NewNotificationPubSub(brokenRDB(t)))
	if err := pubBroken.SendNotification(user.ID, "notifuser", "actor", "actorname", "follow", "hi", ""); err != nil {
		t.Fatalf("publish error should be swallowed, got %v", err)
	}
}

// ---------------- Upload service: validation/error branches ----------------

func TestUploadService_ErrorBranches(t *testing.T) {
	_, db := SetupTestEnv()
	svc := services.NewUploadService(repositories.NewFileRepository(db))
	noop := func(_ *multipart.FileHeader, _ string) error { return nil }

	// Empty file body -> read error.
	empty := makeFileHeader(t, "empty.png", []byte{})
	if _, err := svc.SaveFile(empty, "owner", models.FileVisibilityPublic, noop); err == nil {
		t.Fatal("expected read error on empty file")
	}

	// saveFn failure path.
	valid := makeFileHeader(t, "pic.png", pngBytes(t))
	failSave := func(_ *multipart.FileHeader, _ string) error { return http.ErrNotSupported }
	if _, err := svc.SaveFile(valid, "owner", models.FileVisibilityPublic, failSave); err == nil {
		t.Fatal("expected error when saveFn fails")
	}

	// DB Create failure path (saveFn succeeds, repo is broken).
	bsvc := services.NewUploadService(repositories.NewFileRepository(brokenDB(t)))
	valid2 := makeFileHeader(t, "pic.png", pngBytes(t))
	if _, err := bsvc.SaveFile(valid2, "owner", models.FileVisibilityPublic, noop); err == nil {
		t.Fatal("expected error tracking file in broken db")
	}
}

// ---------------- Friend service: real-DB negative cases ----------------

func TestFriendService_NegativeCases(t *testing.T) {
	_, db := SetupTestEnv()
	svc := &services.FriendService{DB: db}

	// Accept yourself -> error.
	if err := svc.AcceptRequest("u1", "u1"); err == nil {
		t.Fatal("expected error accepting own request")
	}
	// Unfollow when not following -> error.
	if err := svc.Unfollow(utils.NewID(), utils.NewID()); err == nil {
		t.Fatal("expected error unfollowing a non-followed user")
	}
	// Remove a friend that isn't one -> error.
	if err := svc.RemoveFriend(utils.NewID(), utils.NewID()); err == nil {
		t.Fatal("expected error removing a non-friend")
	}
	// Reject a request that doesn't exist -> error.
	if err := svc.RejectRequest(utils.NewID(), utils.NewID()); err == nil {
		t.Fatal("expected error rejecting a missing request")
	}

	// Accept against a broken DB surfaces the non-NotFound error branch.
	bsvc := &services.FriendService{DB: brokenDB(t)}
	if err := bsvc.AcceptRequest("a", "b"); err == nil {
		t.Fatal("expected db error accepting on broken db")
	}
}

// ---------------- User controller: update DB error ----------------

func TestUserController_UpdateDBError(t *testing.T) {
	bcfg, _ := config.Load()
	SetupTestEnv()
	ctrl := routes.Wire(brokenDB(t), sharedRDB, bcfg)

	c, w := testCtx(http.MethodPut, "/", `{"display_name":"new name"}`)
	c.Set("user_id", "u1")
	setParam(c, "id", "u1")
	ctrl.User.UpdateUser(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("UpdateUser db error: expected 500, got %d", w.Code)
	}
}

// ---------------- OAuth controller: redis-error branches ----------------

func TestOAuthController_RedisErrors(t *testing.T) {
	cfg, _ := config.Load()
	_, db := SetupTestEnv()
	ctrl := routes.Wire(db, brokenRDB(t), cfg)

	// Login: GenerateState writes to redis and fails.
	c, w := testCtx(http.MethodGet, "/", "")
	ctrl.OAuth.OAuthLogin(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("OAuthLogin redis error: expected 500, got %d", w.Code)
	}

	// Callback: VerifyAndConsumeState gets a non-Nil redis error.
	c, w = testCtx(http.MethodGet, "/?code=abc&state=xyz", "")
	ctrl.OAuth.OAuthCallback(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("OAuthCallback redis error: expected 500, got %d", w.Code)
	}
}

// ---------------- JWT: remaining branches ----------------

func TestJWT_MoreBranches(t *testing.T) {
	// A token signed with an unexpected (non-HMAC) method is rejected.
	if _, err := utils.ValidateJWT("eyJhbGciOiJub25lIiwidHlwIjoiSldUIn0.eyJzdWIiOiJ4In0."); err == nil {
		t.Fatal("expected error for non-HMAC token")
	}
	// A structurally invalid token is rejected.
	if _, err := utils.ValidateJWT("a.b.c"); err == nil {
		t.Fatal("expected error for garbage token")
	}
	// RefreshToken on an invalid token errors.
	if _, err := utils.RefreshToken("not-a-token"); err == nil {
		t.Fatal("expected refresh error")
	}
	// RefreshToken on a fresh (long-lived) token returns it unchanged.
	tok, err := utils.GenerateJWT(utils.NewID(), "refresher")
	if err != nil {
		t.Fatalf("generate: %v", err)
	}
	got, err := utils.RefreshToken(tok)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if got != tok {
		t.Fatal("expected a long-lived token to be returned unchanged")
	}
}

// ---------------- CreateComment with a file attachment (e2e) ----------------

func TestCreateComment_WithFile(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cmtfile", "cmtfile@test.com", "StrongPass123!")
	postID := createPost(t, router, user.Token, "post for file comment")

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	_ = mw.WriteField("content", "comment with an image")
	_ = mw.WriteField("visibility", models.FileVisibilityPublic)
	part, err := mw.CreateFormFile("file", "pic.png")
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	_, _ = part.Write(pngBytes(t))
	mw.Close()

	req, _ := http.NewRequest(http.MethodPost, "/api/posts/"+postID+"/comments", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment with file: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
}
