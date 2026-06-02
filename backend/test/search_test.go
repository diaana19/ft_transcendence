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

// TestSearch_SortPosts asserts that ?sort=likes_count&order=desc actually
// orders the returned posts — previous tests only checked counts.
func TestSearch_SortPosts(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "sortpost", "sortpost@test.com", "StrongPass123!")
	// Three posts; we will like #2 once and #3 twice so the desc order by
	// likes_count is: post3 (2) -> post2 (1) -> post1 (0).
	id1 := createPost(t, router, author.Token, "sortme alpha")
	id2 := createPost(t, router, author.Token, "sortme beta")
	id3 := createPost(t, router, author.Token, "sortme gamma")

	liker1 := registerAndLogin(t, router, "spliker1", "spliker1@test.com", "StrongPass123!")
	liker2 := registerAndLogin(t, router, "spliker2", "spliker2@test.com", "StrongPass123!")
	if w := authedRequest(t, router, "POST", "/api/posts/"+id2+"/like", liker1.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("like post2: %d", w.Code)
	}
	if w := authedRequest(t, router, "POST", "/api/posts/"+id3+"/like", liker1.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("like post3 a: %d", w.Code)
	}
	if w := authedRequest(t, router, "POST", "/api/posts/"+id3+"/like", liker2.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("like post3 b: %d", w.Code)
	}

	w := authedRequest(t, router, "GET",
		"/api/search?q=sortme&type=post&sort=likes_count&order=desc", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("sort posts: %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Posts) != 3 {
		t.Fatalf("expected 3 posts, got %d", len(res.Posts))
	}
	got := []string{res.Posts[0]["id"].(string), res.Posts[1]["id"].(string), res.Posts[2]["id"].(string)}
	want := []string{id3, id2, id1}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("order by likes_count desc: position %d got %s, want %s (full order %v)", i, got[i], want[i], got)
		}
	}
}

// TestSearch_LikeWildcardLiteral confirms `%` and `_` in the query are matched
// literally — i.e., the SQL LIKE metacharacters are escaped.
func TestSearch_LikeWildcardLiteral(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "wcauthor", "wcauthor@test.com", "StrongPass123!")
	// 'token1' and 'token2' both contain "token", so an UNescaped search for
	// "tok_n" (LIKE: tok + any-char + n) would match both. A correctly
	// escaped query must match zero.
	createPost(t, router, author.Token, "literal token1 here")
	createPost(t, router, author.Token, "literal token2 here")

	w := authedRequest(t, router, "GET", "/api/search?q=tok_n&type=post", author.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search: %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if res.Total != 0 || len(res.Posts) != 0 {
		t.Fatalf("`_` should be matched literally: got total=%d posts=%d", res.Total, len(res.Posts))
	}

	// Likewise `%` should not be a wildcard.
	w = authedRequest(t, router, "GET", "/api/search?q=token%25&type=post", author.Token, "")
	json.Unmarshal(w.Body.Bytes(), &res)
	if res.Total != 0 {
		t.Fatalf("`%%` should be matched literally: got total=%d", res.Total)
	}
}

// TestSearch_PaginationSecondPage asserts that page=2 returns the next slice.
func TestSearch_PaginationSecondPage(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "p_auth", "p_auth@test.com", "StrongPass123!")
	// Five posts. limit=2 -> page1: 2 items, page2: 2 items, page3: 1 item.
	for i := range 5 {
		createPost(t, router, author.Token, "paginate me item "+string(rune('A'+i)))
	}

	w := authedRequest(t, router, "GET", "/api/search?q=paginate&type=post&limit=2&page=1", author.Token, "")
	var p1 searchResponse
	json.Unmarshal(w.Body.Bytes(), &p1)
	if len(p1.Posts) != 2 || p1.Total != 5 {
		t.Fatalf("page1: expected 2 posts of 5 total, got %d/%d", len(p1.Posts), p1.Total)
	}

	w = authedRequest(t, router, "GET", "/api/search?q=paginate&type=post&limit=2&page=2", author.Token, "")
	var p2 searchResponse
	json.Unmarshal(w.Body.Bytes(), &p2)
	if len(p2.Posts) != 2 {
		t.Fatalf("page2: expected 2 posts, got %d", len(p2.Posts))
	}

	w = authedRequest(t, router, "GET", "/api/search?q=paginate&type=post&limit=2&page=3", author.Token, "")
	var p3 searchResponse
	json.Unmarshal(w.Body.Bytes(), &p3)
	if len(p3.Posts) != 1 {
		t.Fatalf("page3: expected 1 post (tail), got %d", len(p3.Posts))
	}

	// page1 ids must be disjoint from page2 ids.
	for _, a := range p1.Posts {
		for _, b := range p2.Posts {
			if a["id"] == b["id"] {
				t.Fatalf("page1 and page2 share id %v — offset not applied", a["id"])
			}
		}
	}
}

// TestSearch_AllIncludesMessages — the "all" branch must search all three
// types and report a coherent total.
func TestSearch_AllIncludesMessages(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "all_alice", "all_alice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "all_bob_omnitag", "all_bob_omnitag@test.com", "StrongPass123!")
	createPost(t, router, alice.Token, "post mentioning omnitag in content")
	body := `{"recipient_id":"` + bob.ID + `","content":"hey omnitag drop by"}`
	if w := authedRequest(t, router, "POST", "/api/chat/messages", alice.Token, body); w.Code != http.StatusCreated {
		t.Fatalf("send message: %d", w.Code)
	}

	w := authedRequest(t, router, "GET", "/api/search?q=omnitag&type=all", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("search all: %d - body: %s", w.Code, w.Body.String())
	}
	var res searchResponse
	json.Unmarshal(w.Body.Bytes(), &res)
	if len(res.Users) != 1 {
		t.Fatalf("expected 1 user, got %d", len(res.Users))
	}
	if len(res.Posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(res.Posts))
	}
	if len(res.Messages) != 1 {
		t.Fatalf("expected 1 message, got %d", len(res.Messages))
	}
	if res.Total != 3 {
		t.Fatalf("expected total=3 (1 user + 1 post + 1 message), got %d", res.Total)
	}
}
