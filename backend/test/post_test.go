package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// createPost creates a post as the given user and returns its id.
func createPost(t *testing.T, router *gin.Engine, token, content string) string {
	t.Helper()
	w := authedRequest(t, router, "POST", "/api/posts", token, `{"content":"`+content+`"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("create post: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		ID      string `json:"id"`
		Content string `json:"content"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode created post: %v", err)
	}
	if resp.ID == "" || resp.Content != content {
		t.Fatalf("unexpected created post: %+v", resp)
	}
	return resp.ID
}

func TestPost_CreateAndGet(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pauthor", "pauthor@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "hello world")

	w := authedRequest(t, router, "GET", "/api/posts/"+id, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get post: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/posts", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list posts: expected 200, got %d", w.Code)
	}
	var list struct {
		Data  []map[string]any `json:"data"`
		Total int64            `json:"total"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &list); err != nil {
		t.Fatalf("decode list: %v", err)
	}
	if list.Total != 1 || len(list.Data) != 1 {
		t.Fatalf("expected 1 post, got total=%d len=%d", list.Total, len(list.Data))
	}
}

func TestPost_GetNotFound(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "pnf", "pnf@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/posts/550e8400-e29b-41d4-a716-446655440000", u.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestPost_CreateRequiresContent(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "pempty", "pempty@test.com", "StrongPass123!")

	w := authedRequest(t, router, "POST", "/api/posts", u.Token, `{}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for missing content, got %d", w.Code)
	}
}

func TestPost_CreateRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "POST", "/api/posts", "", `{"content":"x"}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestPost_UpdateOwnership(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pup", "pup@test.com", "StrongPass123!")
	other := registerAndLogin(t, router, "pother", "pother@test.com", "StrongPass123!")

	id := createPost(t, router, author.Token, "original")

	w := authedRequest(t, router, "PUT", "/api/posts/"+id, author.Token, `{"content":"edited"}`)
	if w.Code != http.StatusOK {
		t.Fatalf("owner update: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "PUT", "/api/posts/"+id, other.Token, `{"content":"hijack"}`)
	if w.Code != http.StatusForbidden {
		t.Fatalf("non-owner update: expected 403, got %d", w.Code)
	}
}

func TestPost_Delete(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pdel", "pdel@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "to delete")

	w := authedRequest(t, router, "DELETE", "/api/posts/"+id, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+id, author.Token, "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("get deleted: expected 404, got %d", w.Code)
	}
}

func TestPost_ToggleLike(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "plike", "plike@test.com", "StrongPass123!")
	liker := registerAndLogin(t, router, "pliker", "pliker@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "like me")

	w := authedRequest(t, router, "POST", "/api/posts/"+id+"/like", liker.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("like: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Liked      bool `json:"liked"`
		LikesCount int  `json:"likes_count"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if !resp.Liked || resp.LikesCount != 1 {
		t.Fatalf("expected liked=true count=1, got %+v", resp)
	}

	w = authedRequest(t, router, "POST", "/api/posts/"+id+"/like", liker.Token, "")
	json.Unmarshal(w.Body.Bytes(), &resp)
	if resp.Liked || resp.LikesCount != 0 {
		t.Fatalf("expected liked=false count=0 after second toggle, got %+v", resp)
	}
}

func TestPost_Comments(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pcom", "pcom@test.com", "StrongPass123!")
	commenter := registerAndLogin(t, router, "pcommenter", "pcommenter@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "comment on me")

	w := authedRequest(t, router, "POST", "/api/posts/"+id+"/comments", commenter.Token, `{"content":"nice post"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment: expected 201, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+id+"/comments", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("list comments: expected 200, got %d", w.Code)
	}
	var list struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &list)
	if list.Total != 1 {
		t.Fatalf("expected 1 comment, got %d", list.Total)
	}
}

func TestPost_GetByUser(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pbyuser", "pbyuser@test.com", "StrongPass123!")
	createPost(t, router, author.Token, "post one")
	createPost(t, router, author.Token, "post two")

	w := authedRequest(t, router, "GET", "/api/posts/user/"+author.ID, author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("get by user: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []map[string]any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 posts by author, got %d", len(resp.Data))
	}
}
