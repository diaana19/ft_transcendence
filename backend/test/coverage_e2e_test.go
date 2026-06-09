package test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func loginAttempt(router *gin.Engine, email, password string) *httptest.ResponseRecorder {
	body := `{"email":"` + email + `","password":"` + password + `"}`
	req, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestForgotPassword_UnknownEmail(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{"email":"nobody@test.com"}`
	req, _ := http.NewRequest("POST", "/api/auth/forgot-password", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("unknown email: expected 200 generic, got %d - %s", w.Code, w.Body.String())
	}
}

func TestForgotPassword_InvalidJSON(t *testing.T) {
	router, _ := SetupTestEnv()

	req, _ := http.NewRequest("POST", "/api/auth/forgot-password", bytes.NewBufferString(`{"email":`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid json: expected 400, got %d", w.Code)
	}
}

func TestPasswordReset_Flow(t *testing.T) {
	router, db := SetupTestEnv()
	u := registerAndLogin(t, router, "pwreset", "pwreset@test.com", "StrongPass123!")

	authSvc := services.NewAuthService(repositories.NewUserRepository(db))

	user, err := authSvc.GetUserByEmail("pwreset@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}
	if user.ID != u.ID {
		t.Fatalf("GetUserByEmail returned wrong user: %s != %s", user.ID, u.ID)
	}

	var delivered string
	if err := authSvc.ResetPassword(user.ID, func(np string) error {
		delivered = np
		return nil
	}); err != nil {
		t.Fatalf("ResetPassword: %v", err)
	}
	if len(delivered) != 16 {
		t.Fatalf("expected 16-char generated password, got %d", len(delivered))
	}

	if w := loginAttempt(router, "pwreset@test.com", "StrongPass123!"); w.Code == http.StatusOK {
		t.Fatal("old password should no longer authenticate")
	}
	if w := loginAttempt(router, "pwreset@test.com", delivered); w.Code != http.StatusOK {
		t.Fatalf("new password login: expected 200, got %d - %s", w.Code, w.Body.String())
	}
}

func TestPasswordReset_DeliverFailure(t *testing.T) {
	router, db := SetupTestEnv()
	registerAndLogin(t, router, "pwfail", "pwfail@test.com", "StrongPass123!")

	authSvc := services.NewAuthService(repositories.NewUserRepository(db))
	user, err := authSvc.GetUserByEmail("pwfail@test.com")
	if err != nil {
		t.Fatalf("GetUserByEmail: %v", err)
	}

	err = authSvc.ResetPassword(user.ID, func(string) error {
		return http.ErrAbortHandler
	})
	if err == nil {
		t.Fatal("expected ResetPassword to fail when delivery fails")
	}

	if w := loginAttempt(router, "pwfail@test.com", "StrongPass123!"); w.Code != http.StatusOK {
		t.Fatalf("password must be unchanged after failed delivery, login got %d", w.Code)
	}
}

func TestPosts_RepliedBy(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "rb-auth", "rb-auth@test.com", "StrongPass123!")
	commenter := registerAndLogin(t, router, "rb-comm", "rb-comm@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "a post worth replying to")
	createComment(t, router, commenter.Token, postID, "my reply")

	w := authedRequest(t, router, "GET", "/api/posts?repliedBy="+commenter.ID, commenter.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("repliedBy list: expected 200, got %d - %s", w.Code, w.Body.String())
	}
	var list struct {
		Data  []map[string]any `json:"data"`
		Total int64            `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode repliedBy list: %v", err)
	}
	if list.Total != 1 || len(list.Data) != 1 {
		t.Fatalf("expected 1 replied post, got total=%d len=%d", list.Total, len(list.Data))
	}

	w = authedRequest(t, router, "GET", "/api/posts?repliedBy="+author.ID, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("repliedBy empty: expected 200, got %d", w.Code)
	}
	json.Unmarshal(w.Body.Bytes(), &list)
	if list.Total != 0 {
		t.Fatalf("author replied to nothing, expected total=0, got %d", list.Total)
	}
}

func updateCommentMultipart(router *gin.Engine, token, postID, commentID, content string, removeFile bool) *httptest.ResponseRecorder {
	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	_ = mw.WriteField("content", content)
	if removeFile {
		_ = mw.WriteField("remove_file", "true")
	}
	mw.Close()

	req, _ := http.NewRequest("PUT", "/api/posts/"+postID+"/comments/"+commentID, &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

func TestComment_UpdateMultipart(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cm-auth", "cm-auth@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post for multipart comment edit")
	commentID := createComment(t, router, author.Token, postID, "first version")

	w := updateCommentMultipart(router, author.Token, postID, commentID, "edited via multipart", true)
	if w.Code != http.StatusOK {
		t.Fatalf("multipart update: expected 200, got %d - %s", w.Code, w.Body.String())
	}
	var resp struct {
		Content string `json:"content"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Content != "edited via multipart" {
		t.Fatalf("expected updated content, got %q", resp.Content)
	}

	w = updateCommentMultipart(router, author.Token, postID, commentID, "", false)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("multipart empty content: expected 400, got %d", w.Code)
	}

	long := make([]byte, 281)
	for i := range long {
		long[i] = 'a'
	}
	w = updateCommentMultipart(router, author.Token, postID, commentID, string(long), false)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("multipart too-long content: expected 400, got %d", w.Code)
	}
}

func TestPostService_CreateValidation(t *testing.T) {
	_, db := SetupTestEnv()
	svc := services.NewPostService(repositories.NewPostRepository(db))

	if _, err := svc.CreatePost("", utils.NewID(), nil, nil); err == nil {
		t.Fatal("empty content should error")
	}
	if _, err := svc.CreatePost(strings.Repeat("a", 281), utils.NewID(), nil, nil); err == nil {
		t.Fatal("too-long content should error")
	}

	if _, total, err := svc.GetRepliedPosts(utils.NewID(), 10, 0); err != nil || total != 0 {
		t.Fatalf("GetRepliedPosts for unknown user: err=%v total=%d", err, total)
	}

	if _, _, err := svc.ReactToPost(utils.NewID(), utils.NewID(), 1); err == nil {
		t.Fatal("ReactToPost on missing post should error")
	}
	if _, _, err := svc.ReactToComment(utils.NewID(), utils.NewID(), 1); err == nil {
		t.Fatal("ReactToComment on missing comment should error")
	}
	if _, err := svc.GetComments(utils.NewID()); err == nil {
		t.Fatal("GetComments on missing post should error")
	}
}

func TestPostReaction_Transitions(t *testing.T) {
	router, db := SetupTestEnv()
	u := registerAndLogin(t, router, "rx-user", "rx-user@test.com", "StrongPass123!")
	psvc := services.NewPostService(repositories.NewPostRepository(db))

	post, err := psvc.CreatePost("reaction transitions target", u.ID, nil, nil)
	if err != nil {
		t.Fatalf("CreatePost: %v", err)
	}

	if v, _, err := psvc.ReactToPost(u.ID, post.ID, 1); err != nil || v != 1 {
		t.Fatalf("like: v=%d err=%v", v, err)
	}
	if v, _, err := psvc.ReactToPost(u.ID, post.ID, -1); err != nil || v != -1 {
		t.Fatalf("swap to dislike: v=%d err=%v", v, err)
	}
	if v, _, err := psvc.ReactToPost(u.ID, post.ID, -1); err != nil || v != 0 {
		t.Fatalf("toggle off: v=%d err=%v", v, err)
	}

	commentID := createComment(t, router, u.Token, post.ID, "a comment to react on")
	if v, _, err := psvc.ReactToComment(u.ID, commentID, 1); err != nil || v != 1 {
		t.Fatalf("comment like: v=%d err=%v", v, err)
	}
	if v, _, err := psvc.ReactToComment(u.ID, commentID, -1); err != nil || v != -1 {
		t.Fatalf("comment swap: v=%d err=%v", v, err)
	}
	if v, _, err := psvc.ReactToComment(u.ID, commentID, -1); err != nil || v != 0 {
		t.Fatalf("comment toggle off: v=%d err=%v", v, err)
	}
}

func TestComment_CreateContentValidation(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cv-auth", "cv-auth@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post for comment validation")

	w := postCommentForm(t, router, author.Token, postID, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty comment: expected 400, got %d", w.Code)
	}

	long := make([]byte, 281)
	for i := range long {
		long[i] = 'b'
	}
	w = postCommentForm(t, router, author.Token, postID, string(long))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("too-long comment: expected 400, got %d", w.Code)
	}
}
