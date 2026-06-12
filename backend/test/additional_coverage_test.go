package test

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"ft_transcendence/backend/internal/utils"
)

func TestAuth_LoginWithUsername(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "loginuser", "loginuser@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/login", "", `{"email":"loginuser","password":"StrongPass123!"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("login with username: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestAuth_LoginWithEmail(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "loginemail", "loginemail@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/login", "", `{"email":"loginemail@test.com","password":"StrongPass123!"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("login with email: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestAuth_LoginWrongPassword(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "loginwp", "loginwp@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/login", "", `{"email":"loginwp@test.com","password":"wrong"}`)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Fatalf("login wrong password: expected 400 or 401, got %d", w.Code)
	}
}

func TestAuth_LoginNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/login", "", `{"email":"nonexistent@test.com","password":"StrongPass123!"}`)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Fatalf("login non-existent: expected 400 or 401, got %d", w.Code)
	}
}

func TestAuth_RegisterDuplicateUsername(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "dupuser", "dup1@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/register", "", `{"username":"dupuser","email":"dup2@test.com","password":"StrongPass123!","dateOfBirth":"2000-01-01"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("duplicate username: expected 400, got %d", w.Code)
	}
}

func TestAuth_RegisterDuplicateEmail(t *testing.T) {
	router, _ := SetupTestEnv()
	registerAndLogin(t, router, "dupemail1", "dupemail@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/auth/register", "", `{"username":"dupemail2","email":"dupemail@test.com","password":"StrongPass123!","dateOfBirth":"2000-01-01"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("duplicate email: expected 400, got %d", w.Code)
	}
}

func TestAuth_RegisterInvalidEmail(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/auth/register", "", `{"username":"invalidemail","email":"notanemail","password":"StrongPass123!","dateOfBirth":"2000-01-01"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid email: expected 400, got %d", w.Code)
	}
}

func TestPost_GetNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "postne", "postne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/"+utils.NewID(), user.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get non-existent post: expected 404, got %d", w.Code)
	}
}

func TestPost_UpdateNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "postune", "postune@test.com", "StrongPass123!")

	w := authedRequest(t, router, "PUT", "/api/posts/"+utils.NewID(), user.Token, `{"content":"updated"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("update non-existent post: expected 404, got %d", w.Code)
	}
}

func TestPost_DeleteNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "postdne", "postdne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+utils.NewID(), user.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("delete non-existent post: expected 404, got %d", w.Code)
	}
}

func TestPost_ReactNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "reactne", "reactne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts/"+utils.NewID()+"/react", user.Token, `{"value":1}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("react non-existent post: expected 404, got %d", w.Code)
	}
}

func TestComment_GetForNonExistentPost(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cgetne", "cgetne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/"+utils.NewID()+"/comments", user.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get comments non-existent post: expected 404, got %d", w.Code)
	}
}

func TestComment_CreateOnNonExistentPost(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "ccrtne", "ccrtne@test.com", "StrongPass123!")

	w := postCommentForm(t, router, user.Token, utils.NewID(), "comment")
	if w.Code != http.StatusNotFound {
		t.Fatalf("create comment non-existent post: expected 404, got %d", w.Code)
	}
}

func TestUpload_InvalidExtensionExtra(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "uploadie", "uploadie@test.com", "StrongPass123!")

	w := uploadFile(t, router, user.Token, "test.exe", "public", []byte{0x00})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("upload invalid extension: expected 400, got %d", w.Code)
	}
}

func TestUpload_InvalidVisibilityExtra(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "uploadiv", "uploadiv@test.com", "StrongPass123!")

	w := uploadFile(t, router, user.Token, "test.png", "invalid", []byte{0x89, 0x50, 0x4E, 0x47})
	if w.Code != http.StatusBadRequest {
		t.Fatalf("upload invalid visibility: expected 400, got %d", w.Code)
	}
}

func TestUser_GetAllUsers(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "getall", "getall@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get all users: expected 200, got %d", w.Code)
	}
	var resp []map[string]any
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp) < 1 {
		t.Fatalf("expected at least 1 user, got %d", len(resp))
	}
}

func TestUser_UpdatePartial(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "updp", "updp@test.com", "StrongPass123!")

	w := authedRequest(t, router, "PUT", "/api/users/"+user.ID, user.Token, `{"bio":"new bio"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update partial: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestFriend_FollowNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "folne", "folne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+utils.NewID(), user.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("follow non-existent: expected 400, got %d", w.Code)
	}
}

func TestFriend_RemoveNonExistentIsNoop(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "rmne", "rmne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/friends/"+utils.NewID(), user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("remove non-existent: expected 200, got %d", w.Code)
	}
}

func TestNotification_MarkAllRead(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "mar-a", "mar-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "mar-b", "mar-b@test.com", "StrongPass123!")

	authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")

	w := authedRequest(t, router, "PATCH", "/api/notification/read", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("mark all read: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	var resp struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 0 {
		t.Fatalf("expected 0 unread, got %d", resp.Total)
	}
}

func TestNotification_GetEmpty(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "notifempty", "notifempty@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/notification", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get empty notifications: expected 200, got %d", w.Code)
	}
	var resp struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Total != 0 {
		t.Fatalf("expected 0 notifications, got %d", resp.Total)
	}
}

func TestGDPR_ExportRequiresAuthExtra(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/gdpr/export", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("export no auth: expected 401, got %d", w.Code)
	}
}

func TestUpload_FileTooLarge(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "uploadtl", "uploadtl@test.com", "StrongPass123!")

	largeData := make([]byte, 26*1024*1024)
	w := uploadFile(t, router, user.Token, "large.png", "public", largeData)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("upload too large: expected 400, got %d", w.Code)
	}
}

func TestUpload_FriendsVisibility(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "uploadfv", "uploadfv@test.com", "StrongPass123!")

	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}
	w := uploadFile(t, router, user.Token, "friends.png", "friends", pngData)
	if w.Code != http.StatusOK {
		t.Fatalf("upload friends visibility: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUpload_PrivateVisibility(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "uploadpv", "uploadpv@test.com", "StrongPass123!")

	pngData := []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A, 0x00, 0x00, 0x00, 0x0D, 0x49, 0x48, 0x44, 0x52}
	w := uploadFile(t, router, user.Token, "private.png", "private", pngData)
	if w.Code != http.StatusOK {
		t.Fatalf("upload private visibility: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestPost_CreateWithFileUpload(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "postfile", "postfile@test.com", "StrongPass123!")

	uploadResp := uploadFile(t, router, user.Token, "post.png", "public", []byte{0x89, 0x50, 0x4E, 0x47})
	var uploadData struct {
		URL string `json:"url"`
	}
	json.Unmarshal(uploadResp.Body.Bytes(), &uploadData)

	w := authedRequest(t, router, "POST", "/api/posts", user.Token, `{"content":"post with file","media_url":"`+uploadData.URL+`","media_mime":"image/png"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("create post with file: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestComment_CreateWithFileUpload(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cmtfile", "cmtfile@test.com", "StrongPass123!")
	postID := createPost(t, router, user.Token, "post for comment file")

	uploadResp := uploadFile(t, router, user.Token, "comment.png", "public", []byte{0x89, 0x50, 0x4E, 0x47})
	var uploadData struct {
		ID string `json:"id"`
	}
	json.Unmarshal(uploadResp.Body.Bytes(), &uploadData)

	var body bytes.Buffer
	mw := multipart.NewWriter(&body)
	mw.WriteField("content", "comment with file")
	mw.Close()

	req, _ := http.NewRequest("POST", "/api/posts/"+postID+"/comments", &body)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("create comment with file: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestPost_GetSingle(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "getsingle", "getsingle@test.com", "StrongPass123!")
	postID := createPost(t, router, user.Token, "single post")

	w := authedRequest(t, router, "GET", "/api/posts/"+postID, user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get single post: expected 200, got %d", w.Code)
	}
	var resp struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Content != "single post" {
		t.Fatalf("expected 'single post', got %q", resp.Content)
	}
}

func TestPost_GetByUserEndpoint(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "getbyuser", "getbyuser@test.com", "StrongPass123!")
	createPost(t, router, user.Token, "user's post")

	w := authedRequest(t, router, "GET", "/api/posts/user/"+user.ID, user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get by user: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) < 1 {
		t.Fatalf("expected at least 1 post, got %d", len(resp.Data))
	}
}
