package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

type searchResponse struct {
	Users    []map[string]any `json:"users"`
	Messages []map[string]any `json:"messages"`
	Posts    []map[string]any `json:"posts"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
}

func TestSearch_Users(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "findzeta", "findzeta@test.com", "StrongPass123!")
	registerAndLogin(t, router, "findzebra", "findzebra@test.com", "StrongPass123!")
	registerAndLogin(t, router, "unrelated", "unrelated@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/search?q=findz&type=user", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search users: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	if err := json.Unmarshal(w.Body.Bytes(), &res); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(res.Users) != 2 {
		t.Fatalf("expected 2 matching users, got %d", len(res.Users))
	}
}

func TestSearch_DefaultTypeIsUser(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "deftype", "deftype@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/search?q=deftype", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("default type: expected 200, got %d", w.Code)
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res.Users))
	}
}

func TestSearch_Posts(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "postsearcher", "postsearcher@test.com", "StrongPass123!")
	createPost(t, router, author.Token, "golang concurrency is great")
	createPost(t, router, author.Token, "rust ownership model")

	w := authedRequest(t, router, "GET", "/api/search?q=golang&type=post", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search posts: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Posts) != 1 {
		t.Fatalf("expected 1 matching post, got %d", len(res.Posts))
	}
}

func TestSearch_Messages(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "msrch_a", "msrch_a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "msrch_b", "msrch_b@test.com", "StrongPass123!")

	body := `{"recipient_id":"` + bob.ID + `","content":"meet me at the pier tonight"}`
	if w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, body); w.Code != http.StatusCreated {
		t.Fatalf("send message: expected 201, got %d", w.Code)
	}

	w := authedRequest(t, router, "GET", "/api/search?q=pier&type=message", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search messages: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Messages) != 1 {
		t.Fatalf("expected 1 matching message, got %d", len(res.Messages))
	}
}

func TestSearch_All(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "omniquux", "omniquux@test.com", "StrongPass123!")
	createPost(t, router, u.Token, "omniquux writes a post")

	w := authedRequest(t, router, "GET", "/api/search?q=omniquux&type=all", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search all: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Users) != 1 || len(res.Posts) != 1 {
		t.Fatalf("search all: expected 1 user and 1 post, got %d users, %d posts", len(res.Users), len(res.Posts))
	}
}

func TestSearch_MissingQuery(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "noquery", "noquery@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/search?type=user", u.Token, "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("missing q: expected 400, got %d", w.Code)
	}
}

func TestSearch_InvalidType(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "badtype", "badtype@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/search?q=x&type=bogus", u.Token, "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("invalid type: expected 500, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestSearch_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/search?q=x&type=user", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("no token: expected 401, got %d", w.Code)
	}
}

func TestSearch_PaginationClamped(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "clamped", "clamped@test.com", "StrongPass123!")

	// page<1 and limit>100 are clamped to defaults; request must still succeed
	w := authedRequest(t, router, "GET", "/api/search?q=clamped&type=user&page=0&limit=9999", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("clamped pagination: expected 200, got %d", w.Code)
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if res.Page != 1 || res.Limit != 20 {
		t.Fatalf("expected clamped page=1 limit=20, got page=%d limit=%d", res.Page, res.Limit)
	}
}
