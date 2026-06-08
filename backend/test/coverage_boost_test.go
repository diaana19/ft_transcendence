package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuth_MeEndpoint(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "me-a", "me-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/auth/me", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("me: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		User struct {
			ID       string `json:"id"`
			Username string `json:"username"`
		} `json:"user"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.User.ID == "" {
		t.Fatal("expected user ID in response")
	}
}

func TestAuth_MeViaCookie(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "me-b", "me-b@test.com", "StrongPass123!")

	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	req.AddCookie(&http.Cookie{Name: "auth_token", Value: user.Token})
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("me via cookie: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestAuth_MeNoAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/auth/me", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("me no auth: expected 401, got %d", w.Code)
	}
}

func TestGDPR_ExportData(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "gdpr-a", "gdpr-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/gdpr/export", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("export: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestGDPR_DeleteData(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "gdpr-b", "gdpr-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/gdpr/delete", user.Token, `{"password":"StrongPass123!"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestGDPR_DeleteNoAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "DELETE", "/api/gdpr/delete", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("delete no auth: expected 401, got %d", w.Code)
	}
}

func TestGamification_GetStats(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "gam-a", "gam-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+user.ID+"/gamification", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("gamification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestGamification_Leaderboard(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "gam-b", "gam-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+user.ID+"/gamification", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("leaderboard: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestNotification_FollowGeneratesNotification(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "notif-follow-a", "notif-follow-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "notif-follow-b", "notif-follow-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("follow: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/notification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []struct {
			Type string `json:"type"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	found := false
	for _, n := range resp.Data {
		if n.Type == "follow" {
			found = true
			break
		}
	}
	if !found {
		t.Fatal("expected follow notification")
	}
}

func TestPost_FilterByTag(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "tag-a", "tag-a@test.com", "StrongPass123!")

	createPost(t, router, user.Token, "post with #testing tag")
	createPost(t, router, user.Token, "another post")

	w := authedRequest(t, router, "GET", "/api/posts?tag=testing", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("filter by tag: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) < 1 {
		t.Fatal("expected at least 1 post with tag")
	}
}

func TestPost_CreateWithMedia(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "media-a", "media-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts", user.Token, `{"content":"post with media","media_url":"/api/files/test","media_mime":"image/png"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("create with media: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		MediaURL  *string `json:"media_url"`
		MediaMIME *string `json:"media_mime"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.MediaURL == nil {
		t.Fatal("expected media_url in response")
	}
}

func TestPost_UpdateWithMedia(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "media-b", "media-b@test.com", "StrongPass123!")

	id := createPost(t, router, user.Token, "original")

	w := authedRequest(t, router, "PUT", "/api/posts/"+id, user.Token, `{"content":"updated","media_url":"/api/files/new","media_mime":"video/mp4"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update with media: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUpload_ImageFile(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "upload-a", "upload-a@test.com", "StrongPass123!")

	w := uploadFile(t, router, user.Token, "test.png", "public", []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A})
	if w.Code != http.StatusOK {
		t.Fatalf("upload image: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUpload_NoFile(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "upload-b", "upload-b@test.com", "StrongPass123!")

	req, _ := http.NewRequest("POST", "/api/upload", &bytes.Buffer{})
	req.Header.Set("Authorization", "Bearer "+user.Token)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("upload no file: expected 400, got %d", w.Code)
	}
}

func TestUser_DeleteAccount(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "del-acc", "del-acc@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/users/"+user.ID, user.Token, `{"password":"StrongPass123!"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("delete account: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUser_DeleteWrongPassword(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "del-wp", "del-wp@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/users/"+user.ID, user.Token, `{"password":"wrong"}`)
	if w.Code != http.StatusBadRequest && w.Code != http.StatusUnauthorized {
		t.Fatalf("delete wrong password: expected 400 or 401, got %d", w.Code)
	}
}

func Test2FA_SetupRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/2fa/setup", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("2fa setup no auth: expected 401, got %d", w.Code)
	}
}

func Test2FA_DisableRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/2fa/disable", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("2fa disable no auth: expected 401, got %d", w.Code)
	}
}

func TestComment_Update(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cmt-upd", "cmt-upd@test.com", "StrongPass123!")

	postID := createPost(t, router, user.Token, "post")
	commentID := createComment(t, router, user.Token, postID, "original comment")

	w := authedRequest(t, router, "PUT", "/api/posts/"+postID+"/comments/"+commentID, user.Token, `{"content":"updated comment"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update comment: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestComment_Delete(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "cmt-del", "cmt-del@test.com", "StrongPass123!")

	postID := createPost(t, router, user.Token, "post")
	commentID := createComment(t, router, user.Token, postID, "to be deleted")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+postID+"/comments/"+commentID, user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete comment: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestComment_DeleteOtherUserFails(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cmt-dof-a", "cmt-dof-a@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "cmt-dof-b", "cmt-dof-b@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "post")
	commentID := createComment(t, router, author.Token, postID, "not yours")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+postID+"/comments/"+commentID, other.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("delete other's comment: expected 403, got %d", w.Code)
	}
}

func TestComment_Like(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cmt-like-a", "cmt-like-a@test.com", "StrongPass123!")
	liker := registerAndLogin(t, router, "cmt-like-b", "cmt-like-b@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "post")
	commentID := createComment(t, router, author.Token, postID, "likable comment")

	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/comments/"+commentID+"/react", liker.Token, `{"value":1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("like comment: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestPost_CreateEmptyContent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "post-empty", "post-empty@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts", user.Token, `{"content":""}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty content: expected 400, got %d", w.Code)
	}
}

func TestPost_CreateNoAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/posts", "", `{"content":"test"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no auth: expected 401, got %d", w.Code)
	}
}

func TestFriend_UnfollowNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "unf-ne", "unf-ne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "DELETE", "/api/friends/follow/550e8400-e29b-41d4-a716-446655440000", user.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unfollow non-existent: expected 400, got %d", w.Code)
	}
}

func TestFriend_RequestSelfFails(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "req-self", "req-self@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/friends/request/"+user.ID, user.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("request self: expected 400, got %d", w.Code)
	}
}

func TestUser_GetNonExistent(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "get-ne", "get-ne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/550e8400-e29b-41d4-a716-446655440000", user.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get non-existent: expected 404, got %d", w.Code)
	}
}

func TestPost_GetByNonExistentUser(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "pbu-ne", "pbu-ne@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/user/550e8400-e29b-41d4-a716-446655440000", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get by non-existent user: expected 200, got %d", w.Code)
	}
}
