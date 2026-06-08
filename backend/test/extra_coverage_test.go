package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestPost_UpdatePost(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "upd-a", "upd-a@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "original content")

	w := authedRequest(t, router, "PUT", "/api/posts/"+id, author.Token, `{"content":"updated content"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update post: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Content string `json:"content"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Content != "updated content" {
		t.Fatalf("expected updated content, got %q", resp.Content)
	}
}

func TestPost_UpdateOtherUserPostFails(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "upd-b", "upd-b@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "upd-c", "upd-c@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "not yours")

	w := authedRequest(t, router, "PUT", "/api/posts/"+id, other.Token, `{"content":"hacked"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("update other's post: expected 403, got %d", w.Code)
	}
}

func TestPost_DeletePost(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "del-a", "del-a@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "to be deleted")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+id, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete post: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+id, author.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get deleted post: expected 404, got %d", w.Code)
	}
}

func TestPost_DeleteOtherUserPostFails(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "del-b", "del-b@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "del-c", "del-c@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "not yours to delete")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+id, other.Token, "")
	if w.Code != http.StatusForbidden {
		t.Fatalf("delete other's post: expected 403, got %d", w.Code)
	}
}

func TestPost_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "pnf-a", "pnf-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/550e8400-e29b-41d4-a716-446655440000", user.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get non-existent post: expected 404, got %d", w.Code)
	}
}

func TestComment_CreateAndGet(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cmt-a", "cmt-a@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "post with comments")

	// Comments use multipart form data
	w := postCommentForm(t, router, author.Token, postID, "nice post")
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+postID+"/comments", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get comments: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 1 {
		t.Fatalf("expected 1 comment, got %d", len(resp.Data))
	}
}

func TestComment_EmptyContentFails(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "cmt-b", "cmt-b@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "post")

	w := postCommentForm(t, router, author.Token, postID, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("empty comment: expected 400, got %d", w.Code)
	}
}

func TestUser_GetUser(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "usr-a", "usr-a@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+user.ID, user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get user: expected 200, got %d", w.Code)
	}
}

func TestUser_GetUsers(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "usr-b", "usr-b@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users", user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get users: expected 200, got %d", w.Code)
	}
}

func TestUser_UpdateUser(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "usr-c", "usr-c@test.com", "StrongPass123!")

	w := authedRequest(t, router, "PUT", "/api/users/"+user.ID, user.Token, `{"displayname":"New Name","bio":"New bio"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("update user: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestUser_UpdateOtherUserFails(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "usr-d", "usr-d@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "usr-e", "usr-e@test.com", "StrongPass123!")

	w := authedRequest(t, router, "PUT", "/api/users/"+user.ID, other.Token, `{"displayname":"Hacked"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("update other user: expected 403, got %d", w.Code)
	}
}

func TestTrends_ReturnsData(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "trend-a", "trend-a@test.com", "StrongPass123!")

	createPost(t, router, user.Token, "post with #golang tag")

	w := authedRequest(t, router, "GET", "/api/trends", "", "")
	if w.Code != http.StatusOK {
		t.Fatalf("trends: expected 200, got %d", w.Code)
	}
}

func TestPost_GetByUserExtra(t *testing.T) {
	router, _ := SetupTestEnv()
	user := registerAndLogin(t, router, "pbu-a", "pbu-a@test.com", "StrongPass123!")

	createPost(t, router, user.Token, "user's post")

	w := authedRequest(t, router, "GET", "/api/posts/user/"+user.ID, user.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("posts by user: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) < 1 {
		t.Fatalf("expected at least 1 post, got %d", len(resp.Data))
	}
}

func TestPost_LikeAndUnlike(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "react-a", "react-a@test.com", "StrongPass123!")
	liker := registerAndLogin(t, router, "react-b", "react-b@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "likeable post")

	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/react", liker.Token, `{"value":1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("like: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	// Unlike (toggle by sending same value again)
	w = authedRequest(t, router, "POST", "/api/posts/"+postID+"/react", liker.Token, `{"value":1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("unlike: expected 200, got %d", w.Code)
	}
}

func TestPost_Dislike(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "react-c", "react-c@test.com", "StrongPass123!")
	disliker := registerAndLogin(t, router, "react-d", "react-d@test.com", "StrongPass123!")

	postID := createPost(t, router, author.Token, "dislikeable post")

	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/react", disliker.Token, `{"value":-1}`)
	if w.Code != http.StatusOK {
		t.Fatalf("dislike: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestPost_InvalidReactionValue(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "react-e", "react-e@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "post")

	w := authedRequest(t, router, "POST", "/api/posts/"+postID+"/react", author.Token, `{"value":5}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid reaction: expected 400, got %d", w.Code)
	}
}

func TestPost_Pagination(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "page-a", "page-a@test.com", "StrongPass123!")

	for i := 0; i < 5; i++ {
		createPost(t, router, author.Token, "post "+string(rune('a'+i)))
	}

	w := authedRequest(t, router, "GET", "/api/posts?limit=2&offset=0", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("page 1: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data  []map[string]any `json:"data"`
		Total int64            `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(resp.Data))
	}

	w = authedRequest(t, router, "GET", "/api/posts?limit=2&offset=2", author.Token, "")
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 posts on page 2, got %d", len(resp.Data))
	}
}
