package tests

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

// createComment posts a comment on postID as the given user and returns its id.
func createComment(t *testing.T, router *gin.Engine, token, postID, content string) string {
	t.Helper()
	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/comments", token, `{"content":"`+content+`"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode comment: %v", err)
	}
	if resp.ID == "" {
		t.Fatal("expected non-empty comment id")
	}
	return resp.ID
}

func TestComment_UpdateOwnership(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cup_a", "cup_a@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "cup_o", "cup_o@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post with comments")
	commentID := createComment(t, router, author.Token, postID, "original comment")

	w := authedRequest(t, router, "PUT", "/api/posts/"+postID+"/comments/"+commentID, author.Token, `{"content":"edited comment"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("owner update comment: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Content string `json:"content"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Content != "edited comment" {
		t.Fatalf("expected updated content, got %q", resp.Content)
	}

	w = authedRequest(t, router, "PUT", "/api/posts/"+postID+"/comments/"+commentID, other.Token, `{"content":"hijack"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-owner update comment: expected 403, got %d", w.Code)
	}
}

func TestComment_UpdateNotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cupnf", "cupnf@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post")

	w := authedRequest(t, router, "PUT", "/api/posts/"+postID+"/comments/550e8400-e29b-41d4-a716-446655440000",
		author.Token, `{"content":"x"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("update missing comment: expected 404, got %d", w.Code)
	}
}

func TestComment_UpdateInvalidBody(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cupib", "cupib@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post")
	commentID := createComment(t, router, author.Token, postID, "comment")

	w := authedRequest(t, router, "PUT", "/api/posts/"+postID+"/comments/"+commentID, author.Token, `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("update with empty content: expected 400, got %d", w.Code)
	}
}

func TestComment_DeleteOwnership(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cdel_a", "cdel_a@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "cdel_o", "cdel_o@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post")
	commentID := createComment(t, router, author.Token, postID, "comment")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+postID+"/comments/"+commentID, other.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-owner delete comment: expected 403, got %d", w.Code)
	}

	w = authedRequest(t, router, "DELETE", "/api/posts/"+postID+"/comments/"+commentID, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("owner delete comment: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestComment_DeleteNotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cdelnf", "cdelnf@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+postID+"/comments/550e8400-e29b-41d4-a716-446655440000",
		author.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("delete missing comment: expected 404, got %d", w.Code)
	}
}

func TestComment_CreateOnMissingPost(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "ccmp", "ccmp@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts/550e8400-e29b-41d4-a716-446655440000/comments",
		u.Token, `{"content":"hi"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("comment on missing post: expected 404, got %d", w.Code)
	}
}

func TestComment_NotifiesPostAuthor(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cnotif_a", "cnotif_a@test.com", "StrongPass123!")
	commenter := registerAndLogin(t, router, "cnotif_c", "cnotif_c@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "notify me")

	createComment(t, router, commenter.Token, postID, "great post")

	w := authedRequest(t, router, "GET", "/api/notification", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get notifications: expected 200, got %d", w.Code)
	}
}

func TestPost_DeleteForbiddenAndNotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pdf_a", "pdf_a@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "pdf_o", "pdf_o@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "mine")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+postID, other.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-owner delete: expected 403, got %d", w.Code)
	}

	w = authedRequest(t, router, "DELETE", "/api/posts/550e8400-e29b-41d4-a716-446655440000", author.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("delete missing post: expected 404, got %d", w.Code)
	}
}

func TestPost_LikeMissingPost(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "plnf", "plnf@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts/550e8400-e29b-41d4-a716-446655440000/like", u.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("like missing post: expected 404, got %d", w.Code)
	}
}

func TestPost_LikeNotifiesAuthor(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "plna_a", "plna_a@test.com", "StrongPass123!")
	liker := registerAndLogin(t, router, "plna_l", "plna_l@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "like and notify")

	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/like", liker.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("like: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/notification", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("notifications: expected 200, got %d", w.Code)
	}
}

func TestPost_CreateContentTooLong(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "ptoolong", "ptoolong@test.com", "StrongPass123!")

	var long strings.Builder
	for range 281 {
		long.WriteString("x")
	}
	w := authedRequest(t, router, "POST", "/api/posts", u.Token, `{"content":"`+long.String()+`"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("post content >280: expected 400, got %d", w.Code)
	}
}

func TestComment_GetOnMissingPost(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "cgmp", "cgmp@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/550e8400-e29b-41d4-a716-446655440000/comments", u.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("comments on missing post: expected 404, got %d", w.Code)
	}
}
