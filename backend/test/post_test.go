package test

import (
	"encoding/json"
	"net/http"
	"strconv"
	"testing"

	"github.com/gin-gonic/gin"
)

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

func TestPost_ListLimitOffset(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "ppage", "ppage@test.com", "StrongPass123!")
	for _, c := range []string{"one", "two", "three"} {
		createPost(t, router, u.Token, c)
	}

	type listResp struct {
		Data   []map[string]any `json:"data"`
		Limit  int              `json:"limit"`
		Offset int              `json:"offset"`
		Total  int64            `json:"total"`
	}
	get := func(query string) listResp {
		w := authedRequest(t, router, "GET", "/api/posts"+query, u.Token, "")
		if w.Code != http.StatusOK {
			t.Fatalf("list %q: expected 200, got %d", query, w.Code)
		}
		var r listResp
		json.Unmarshal(w.Body.Bytes(), &r)
		return r
	}

	if r := get("?limit=2"); len(r.Data) != 2 || r.Limit != 2 || r.Offset != 0 || r.Total != 3 {
		t.Fatalf("limit=2: got len=%d limit=%d offset=%d total=%d", len(r.Data), r.Limit, r.Offset, r.Total)
	}
	if r := get("?limit=2&offset=2"); len(r.Data) != 1 || r.Offset != 2 {
		t.Fatalf("offset=2: got len=%d offset=%d", len(r.Data), r.Offset)
	}
	if r := get("?limit=9999&offset=-5"); r.Limit != 50 || r.Offset != 0 {
		t.Fatalf("clamp: got limit=%d offset=%d", r.Limit, r.Offset)
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

type reactionResp struct {
	UserReaction  int `json:"user_reaction"`
	LikesCount    int `json:"likes_count"`
	DislikesCount int `json:"dislikes_count"`
}

func react(t *testing.T, router *gin.Engine, path, token string, value int) reactionResp {
	t.Helper()
	w := authedRequest(t, router, "POST", path, token, `{"value":`+strconv.Itoa(value)+`}`)
	if w.Code != http.StatusOK {
		t.Fatalf("react %d on %s: expected 200, got %d - body: %s", value, path, w.Code, w.Body.String())
	}
	var resp reactionResp
	json.Unmarshal(w.Body.Bytes(), &resp)
	return resp
}

func TestPost_LikeDislike(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "plike", "plike@test.com", "StrongPass123!")
	liker := registerAndLogin(t, router, "pliker", "pliker@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "like me")
	path := "/api/posts/" + id + "/react"

	if r := react(t, router, path, liker.Token, 1); r.UserReaction != 1 || r.LikesCount != 1 || r.DislikesCount != 0 {
		t.Fatalf("expected reaction=1 likes=1 dislikes=0, got %+v", r)
	}
	if r := react(t, router, path, liker.Token, -1); r.UserReaction != -1 || r.LikesCount != 0 || r.DislikesCount != 1 {
		t.Fatalf("expected reaction=-1 likes=0 dislikes=1, got %+v", r)
	}
	if r := react(t, router, path, liker.Token, -1); r.UserReaction != 0 || r.LikesCount != 0 || r.DislikesCount != 0 {
		t.Fatalf("expected reaction=0 likes=0 dislikes=0 after toggle off, got %+v", r)
	}
}

func TestPost_ReactValidation(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "prv", "prv@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "react to me")

	for _, body := range []string{`{"value":2}`, `{"value":0}`, `{}`} {
		w := authedRequest(t, router, "POST", "/api/posts/"+id+"/react", author.Token, body)
		if w.Code != http.StatusBadRequest {
			t.Fatalf("invalid reaction body %s: expected 400, got %d", body, w.Code)
		}
	}

	w := authedRequest(t, router, "POST", "/api/posts/550e8400-e29b-41d4-a716-446655440000/react", author.Token, `{"value":1}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("react on missing post: expected 404, got %d", w.Code)
	}
}

func TestPost_ReactRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "prra", "prra@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "react needs auth")

	w := authedRequest(t, router, "POST", "/api/posts/"+id+"/react", "", `{"value":1}`)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("react without token: expected 401, got %d", w.Code)
	}
}

func TestPost_AuthorCanReactOwnPostNoSelfNotify(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "psr", "psr@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "react to myself")

	if r := react(t, router, "/api/posts/"+id+"/react", author.Token, 1); r.UserReaction != 1 || r.LikesCount != 1 {
		t.Fatalf("author self-like: expected reaction=1 likes=1, got %+v", r)
	}

	w := authedRequest(t, router, "GET", "/api/notification", author.Token, "")
	var n struct {
		Total int `json:"total"`
	}
	json.Unmarshal(w.Body.Bytes(), &n)
	if n.Total != 0 {
		t.Fatalf("self-like should not notify: expected 0 notifications, got %d", n.Total)
	}
}

func TestPost_ReactionsAggregateAcrossUsers(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pagg", "pagg@test.com", "StrongPass123!")
	u1 := registerAndLogin(t, router, "pagg1", "pagg1@test.com", "StrongPass123!")
	u2 := registerAndLogin(t, router, "pagg2", "pagg2@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "count me")
	path := "/api/posts/" + id + "/react"

	react(t, router, path, u1.Token, 1)
	react(t, router, path, u2.Token, 1)
	if r := react(t, router, path, u1.Token, -1); r.LikesCount != 1 || r.DislikesCount != 1 {
		t.Fatalf("expected likes=1 dislikes=1 after one switches, got %+v", r)
	}
}

func TestPost_CommentLikeDislike(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pcl", "pcl@test.com", "StrongPass123!")
	reactor := registerAndLogin(t, router, "pclr", "pclr@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "comment then react")

	w := postCommentForm(t, router, author.Token, postID, "a comment")
	if w.Code != http.StatusCreated {
		t.Fatalf("create comment: expected 201, got %d", w.Code)
	}
	var created struct {
		ID string `json:"id"`
	}
	json.Unmarshal(w.Body.Bytes(), &created)

	path := "/api/posts/" + postID + "/comments/" + created.ID + "/react"
	if r := react(t, router, path, reactor.Token, 1); r.UserReaction != 1 || r.LikesCount != 1 {
		t.Fatalf("expected comment reaction=1 likes=1, got %+v", r)
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+postID+"/comments", reactor.Token, "")
	var list struct {
		Data []struct {
			Liked         bool `json:"liked"`
			Disliked      bool `json:"disliked"`
			LikesCount    int  `json:"likes_count"`
			DislikesCount int  `json:"dislikes_count"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list.Data) != 1 || !list.Data[0].Liked || list.Data[0].LikesCount != 1 {
		t.Fatalf("expected liked comment with likes=1, got %+v", list.Data)
	}

	if r := react(t, router, path, reactor.Token, -1); r.UserReaction != -1 || r.LikesCount != 0 || r.DislikesCount != 1 {
		t.Fatalf("expected comment reaction=-1 likes=0 dislikes=1, got %+v", r)
	}
	if r := react(t, router, path, reactor.Token, -1); r.UserReaction != 0 || r.LikesCount != 0 || r.DislikesCount != 0 {
		t.Fatalf("expected comment reaction cleared, got %+v", r)
	}

	w = authedRequest(t, router, "GET", "/api/posts/"+postID+"/comments", "", "")
	json.Unmarshal(w.Body.Bytes(), &list)
	if len(list.Data) != 1 || list.Data[0].Liked || list.Data[0].Disliked {
		t.Fatalf("anon viewer should see no reaction flags, got %+v", list.Data)
	}
}

func TestPost_CommentReactValidation(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pcrv", "pcrv@test.com", "StrongPass123!")
	postID := createPost(t, router, author.Token, "comment validation")

	missing := "/api/posts/" + postID + "/comments/550e8400-e29b-41d4-a716-446655440000/react"
	w := authedRequest(t, router, "POST", missing, author.Token, `{"value":1}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("react on missing comment: expected 404, got %d", w.Code)
	}

	w = authedRequest(t, router, "POST", missing, author.Token, `{"value":5}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("invalid comment reaction value: expected 400, got %d", w.Code)
	}
}

func TestPost_Comments(t *testing.T) {
	router, _ := SetupTestEnv()
	author := registerAndLogin(t, router, "pcom", "pcom@test.com", "StrongPass123!")
	commenter := registerAndLogin(t, router, "pcommenter", "pcommenter@test.com", "StrongPass123!")
	id := createPost(t, router, author.Token, "comment on me")

	w := postCommentForm(t, router, commenter.Token, id, "nice post")
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

func TestPost_Trends(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "ptrend", "ptrend@test.com", "StrongPass123!")

	createPost(t, router, u.Token, "love #golang and #Golang and #redis")
	createPost(t, router, u.Token, "more #golang here")
	createPost(t, router, u.Token, "no tags at all")

	w := authedRequest(t, router, "GET", "/api/trends", "", "")
	if w.Code != http.StatusOK {
		t.Fatalf("trends: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	var resp struct {
		Data []struct {
			Tag   string `json:"tag"`
			Count int64  `json:"count"`
		} `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 2 {
		t.Fatalf("expected 2 distinct tags, got %d: %+v", len(resp.Data), resp.Data)
	}
	if resp.Data[0].Tag != "#golang" || resp.Data[0].Count != 2 {
		t.Fatalf("expected top tag #golang count=2, got %+v", resp.Data[0])
	}
	if resp.Data[1].Tag != "#redis" || resp.Data[1].Count != 1 {
		t.Fatalf("expected #redis count=1, got %+v", resp.Data[1])
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
