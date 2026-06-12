package test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/routes"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/socket"
	"ft_transcendence/backend/internal/utils"
)

func realControllers(t *testing.T) *routes.Controllers {
	t.Helper()
	_, db := SetupTestEnv()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	return routes.Wire(db, sharedRDB, cfg)
}

func TestPostRepository_ReactionStateMachine(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	uid := utils.NewID()
	pid := utils.NewID()
	db.Create(&models.User{ID: uid, Username: "reactsm", Email: "reactsm@test.com"})
	if err := repo.Create(&models.Post{ID: pid, AuthorID: uid, Content: "x"}); err != nil {
		t.Fatalf("create post: %v", err)
	}

	if err := repo.SetPostReaction(uid, pid, 1); err != nil {
		t.Fatalf("set like: %v", err)
	}
	if err := repo.SetPostReaction(uid, pid, 1); err != nil {
		t.Fatalf("set like again: %v", err)
	}
	if err := repo.SetPostReaction(uid, pid, -1); err != nil {
		t.Fatalf("switch to dislike: %v", err)
	}
	if err := repo.SetPostReaction(uid, pid, 0); err != nil {
		t.Fatalf("clear reaction: %v", err)
	}

	reply := &models.Reply{ID: utils.NewID(), PostID: pid, AuthorID: uid, Content: "c"}
	if err := repo.CreateComment(reply); err != nil {
		t.Fatalf("create comment: %v", err)
	}
	rid := reply.ID
	if err := repo.SetReplyReaction(uid, rid, 1); err != nil {
		t.Fatalf("reply like: %v", err)
	}
	if err := repo.SetReplyReaction(uid, rid, 1); err != nil {
		t.Fatalf("reply like again: %v", err)
	}
	if err := repo.SetReplyReaction(uid, rid, -1); err != nil {
		t.Fatalf("reply switch: %v", err)
	}
	if err := repo.SetReplyReaction(uid, rid, 0); err != nil {
		t.Fatalf("reply clear: %v", err)
	}

	if err := repo.Delete(utils.NewID()); err == nil {
		t.Fatal("expected error deleting missing post")
	}
}

func TestPostRepository_DBErrors(t *testing.T) {
	SetupTestEnv()
	repo := repositories.NewPostRepository(brokenDB(t))

	if _, _, err := repo.GetAll(10, 0); err == nil {
		t.Fatal("GetAll: expected error")
	}
	if _, err := repo.GetByID("x"); err == nil {
		t.Fatal("GetByID: expected error")
	}
	if _, err := repo.GetByAuthorID("x"); err == nil {
		t.Fatal("GetByAuthorID: expected error")
	}
	if _, _, err := repo.GetByTag("#t", 10, 0); err == nil {
		t.Fatal("GetByTag: expected error")
	}
	if _, err := repo.TopTags(time.Time{}, 10); err == nil {
		t.Fatal("TopTags: expected error")
	}
	if _, err := repo.Update("x", models.UpdatePostInput{Content: "y"}); err == nil {
		t.Fatal("Update: expected error")
	}
	if err := repo.Delete("x"); err == nil {
		t.Fatal("Delete: expected error")
	}
	if err := repo.CreateComment(&models.Reply{ID: utils.NewID(), PostID: "p", AuthorID: "a", Content: "c"}); err == nil {
		t.Fatal("CreateComment: expected error")
	}
	if _, err := repo.GetCommentsByPostID("x"); err == nil {
		t.Fatal("GetCommentsByPostID: expected error")
	}
	if _, err := repo.GetCommentByID("x"); err == nil {
		t.Fatal("GetCommentByID: expected error")
	}
	if _, err := repo.UpdateComment("x", models.UpdateCommentInput{Content: "y"}); err == nil {
		t.Fatal("UpdateComment: expected error")
	}
	if err := repo.DeleteComment("x"); err == nil {
		t.Fatal("DeleteComment: expected error")
	}
	if _, err := repo.GetPostReaction("u", "p"); err == nil {
		t.Fatal("GetPostReaction: expected error")
	}
	if _, err := repo.GetReplyReaction("u", "r"); err == nil {
		t.Fatal("GetReplyReaction: expected error")
	}
}

func doJSON(router interface {
	ServeHTTP(http.ResponseWriter, *http.Request)
}, method, path, body string,
) *httptest.ResponseRecorder {
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestPostService_EmptyContent(t *testing.T) {
	_, db := SetupTestEnv()
	svc := services.NewPostService(repositories.NewPostRepository(db))
	if _, err := svc.CreatePost("", "author", nil, nil); err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestAuthController_RegisterValidation(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"username":"-bad-","email":"x@test.com","password":"Abcdef1!","dateOfBirth":"2000-01-01"}`
	w := doJSON(router, http.MethodPost, "/api/auth/register", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad username: expected 400, got %d", w.Code)
	}

	body = `{"username":"gooduser","email":"x@test.com","password":"Abcdef1!","dateOfBirth":"not-a-date"}`
	w = doJSON(router, http.MethodPost, "/api/auth/register", body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("bad date: expected 400, got %d", w.Code)
	}
}

func TestAuthController_MeAndLogout(t *testing.T) {
	ctrl := realControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	ctrl.Auth.Me(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Me unauth: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	c.Set("user_id", utils.NewID())
	ctrl.Auth.Me(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Me not found: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", "")
	ctrl.Auth.LogoutUser(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Logout missing token: expected 200, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", "")
	c.Request.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	ctrl.Auth.LogoutUser(c)
	if w.Code != http.StatusOK {
		t.Fatalf("Logout invalid token: expected 200, got %d", w.Code)
	}
}

func TestTwoFAController_Unauthorized(t *testing.T) {
	ctrl := realControllers(t)

	c, w := testCtx(http.MethodPost, "/", "")
	ctrl.TwoFA.Setup(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Setup: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"code":"123456"}`)
	ctrl.TwoFA.Enable(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Enable: expected 401, got %d", w.Code)
	}

	c, w = testCtx(http.MethodPost, "/", `{"code":"123456"}`)
	ctrl.TwoFA.Disable(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("Disable: expected 401, got %d", w.Code)
	}
}

func TestUploadController_Branches(t *testing.T) {
	ctrl := realControllers(t)
	_, db := SetupTestEnv()

	c, w := testCtx(http.MethodPost, "/", "")
	ctrl.Upload.UploadFile(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("UploadFile unauth: expected 401, got %d", w.Code)
	}

	bcfg, _ := config.Load()
	broken := routes.Wire(brokenDB(t), sharedRDB, bcfg)
	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", "anything")
	broken.Upload.ServeFile(c)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("ServeFile invalid id: expected 400, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", utils.NewID())
	broken.Upload.ServeFile(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("ServeFile db error: expected 500, got %d", w.Code)
	}

	fileRepo := repositories.NewFileRepository(db)
	owner := utils.NewID()
	db.Create(&models.User{ID: owner, Username: "upowner", Email: "upowner@test.com"})
	missing := &models.File{
		ID: utils.NewID(), OwnerID: owner, Path: "/tmp/does-not-exist-" + utils.NewID(),
		Filename: "x.png", MimeType: "image/png", Size: 1, Visibility: models.FileVisibilityPublic,
	}
	if err := fileRepo.Create(missing); err != nil {
		t.Fatalf("create file: %v", err)
	}
	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", missing.ID)
	ctrl.Upload.ServeFile(c)
	if w.Code != http.StatusNotFound {
		t.Fatalf("ServeFile missing disk: expected 404, got %d", w.Code)
	}

	tmp, err := os.CreateTemp(t.TempDir(), "f-*.png")
	if err != nil {
		t.Fatalf("temp file: %v", err)
	}
	_, _ = tmp.WriteString("data")
	_ = tmp.Close()
	friendsFile := &models.File{
		ID: utils.NewID(), OwnerID: owner, Path: tmp.Name(),
		Filename: "f.png", MimeType: "image/png", Size: 4, Visibility: models.FileVisibilityFriends,
	}
	if err := fileRepo.Create(friendsFile); err != nil {
		t.Fatalf("create friends file: %v", err)
	}
	c, w = testCtx(http.MethodGet, "/", "")
	c.Set("user_id", utils.NewID())
	setParam(c, "id", friendsFile.ID)
	ctrl.Upload.ServeFile(c)
	if w.Code != http.StatusForbidden {
		t.Fatalf("ServeFile friends forbidden: expected 403, got %d", w.Code)
	}

	privFile := &models.File{
		ID: utils.NewID(), OwnerID: owner, Path: tmp.Name(),
		Filename: "p.png", MimeType: "image/png", Size: 4, Visibility: models.FileVisibilityPrivate,
	}
	if err := fileRepo.Create(privFile); err != nil {
		t.Fatalf("create private file: %v", err)
	}
	c, w = testCtx(http.MethodGet, "/", "")
	setParam(c, "id", privFile.ID)
	ctrl.Upload.ServeFile(c)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("ServeFile private unauth: expected 401, got %d", w.Code)
	}
}

func TestOAuthController_CallbackRedirects(t *testing.T) {
	ctrl := realControllers(t)

	c, w := testCtx(http.MethodGet, "/", "")
	ctrl.OAuth.OAuthCallback(c)
	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("callback no code: expected 307, got %d", w.Code)
	}

	c, w = testCtx(http.MethodGet, "/?code=abc&state=never-stored", "")
	ctrl.OAuth.OAuthCallback(c)
	if w.Code != http.StatusTemporaryRedirect {
		t.Fatalf("callback bad state: expected 307, got %d", w.Code)
	}
}

func TestModels_ToResponseVariants(t *testing.T) {
	mime := "image/png"
	fid := utils.NewID()
	reply := models.Reply{
		ID: utils.NewID(), PostID: "p", AuthorID: "a", Content: "c",
		FileID: &fid, File: &models.File{ID: fid, MimeType: mime},
	}
	resp := reply.ToResponse()
	if resp.FileURL == nil || resp.FileMIME == nil {
		t.Fatal("expected file fields populated on comment response")
	}

	repost := models.Repost{ID: utils.NewID(), PostID: "p", AuthorID: "a"}
	if repost.ToResponse().ID != repost.ID {
		t.Fatal("repost ToResponse mismatch")
	}
}

func TestJWT_ShortSecret(t *testing.T) {
	old := os.Getenv("JWT_SECRET")
	defer os.Setenv("JWT_SECRET", old)

	os.Setenv("JWT_SECRET", "too-short")
	if _, err := utils.GenerateJWT("u", "n"); err == nil {
		t.Fatal("GenerateJWT: expected error with short secret")
	}
	if _, err := utils.ValidateJWT("whatever"); err == nil {
		t.Fatal("ValidateJWT: expected error with short secret")
	}
}

func TestChatHandler_HandleMessage(t *testing.T) {
	ctrl := realControllers(t)
	_, db := SetupTestEnv()
	h := ctrl.ChatWS

	sender := utils.NewID()
	recipient := utils.NewID()
	db.Create(&models.User{ID: sender, Username: "chatsender", Email: "cs@test.com"})
	db.Create(&models.User{ID: recipient, Username: "chatrecipient", Email: "cr@test.com"})

	client := &socket.Client{ID: sender, Username: "chatsender", Send: make(chan []byte, 256)}

	h.HandleMessage(client, []byte("not json"))
	h.HandleMessage(client, []byte(`{"action":"bogus"}`))
	h.HandleMessage(client, []byte(`{"action":"open","peer_id":""}`))
	h.HandleMessage(client, []byte(`{"action":"open","peer_id":"`+sender+`"}`))
	h.HandleMessage(client, []byte(`{"action":"message","content":"   "}`))
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+sender+`"}`))
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+utils.NewID()+`"}`))
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+recipient+`","file_id":"`+utils.NewID()+`"}`))
	h.HandleMessage(client, []byte(`{"action":"open","peer_id":"`+recipient+`"}`))
	h.HandleMessage(client, []byte(`{"action":"message","content":"hello there","recipient_id":"`+recipient+`"}`))

	for {
		select {
		case <-client.Send:
		default:
			return
		}
	}
}

func TestChatHandler_AttachmentValidation(t *testing.T) {
	ctrl := realControllers(t)
	_, db := SetupTestEnv()
	h := ctrl.ChatWS

	sender := utils.NewID()
	recipient := utils.NewID()
	other := utils.NewID()
	db.Create(&models.User{ID: sender, Username: "atsender", Email: "ats@test.com"})
	db.Create(&models.User{ID: recipient, Username: "atrecipient", Email: "atr@test.com"})
	db.Create(&models.User{ID: other, Username: "atother", Email: "ato@test.com"})

	fileRepo := repositories.NewFileRepository(db)
	client := &socket.Client{ID: sender, Username: "atsender", Send: make(chan []byte, 256)}

	notOwned := &models.File{ID: utils.NewID(), OwnerID: other, Path: "/tmp/x", Filename: "x", MimeType: "image/png", Size: 1, Visibility: models.FileVisibilityPrivate}
	fileRepo.Create(notOwned)
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+recipient+`","file_id":"`+notOwned.ID+`"}`))

	pub := &models.File{ID: utils.NewID(), OwnerID: sender, Path: "/tmp/y", Filename: "y", MimeType: "image/png", Size: 1, Visibility: models.FileVisibilityPublic}
	fileRepo.Create(pub)
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+recipient+`","file_id":"`+pub.ID+`"}`))

	priv := &models.File{ID: utils.NewID(), OwnerID: sender, Path: "/tmp/z", Filename: "z", MimeType: "image/png", Size: 1, Visibility: models.FileVisibilityPrivate}
	fileRepo.Create(priv)
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+recipient+`","file_id":"`+priv.ID+`"}`))

	for {
		select {
		case <-client.Send:
		default:
			return
		}
	}
}
