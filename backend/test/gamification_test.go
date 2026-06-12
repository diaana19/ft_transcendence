package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"ft_transcendence/backend/internal/models"
)

type metricResponse struct {
	Level int   `json:"level"`
	Count int64 `json:"count"`
}

type gamificationResponse struct {
	Level     int            `json:"level"`
	Total     int64          `json:"total"`
	Posts     metricResponse `json:"posts"`
	Likes     metricResponse `json:"likes"`
	Followers metricResponse `json:"followers"`
	Following metricResponse `json:"following"`
}

func decodeGamification(t *testing.T, body []byte) gamificationResponse {
	t.Helper()
	var res gamificationResponse
	if err := json.Unmarshal(body, &res); err != nil {
		t.Fatalf("decode gamification response: %v", err)
	}
	return res
}

func TestGamification_NewUserScoresZero(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "gamerzero", "gamerzero@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+u.ID+"/gamification", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	stats := decodeGamification(t, w.Body.Bytes())
	zero := metricResponse{Level: 0, Count: 0}
	if stats.Total != 0 || stats.Level != 0 ||
		stats.Posts != zero || stats.Likes != zero || stats.Followers != zero || stats.Following != zero {
		t.Fatalf("new user should be all zeros, got %+v", stats)
	}
}

func TestGamification_AggregatesPostsLikesFollowersFollowing(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "gameralice", "gameralice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "gamerbob", "gamerbob@test.com", "StrongPass123!")

	post1 := createPost(t, router, alice.Token, "first")
	post2 := createPost(t, router, alice.Token, "second")
	createPost(t, router, alice.Token, "third")

	if w := authedRequest(t, router, "POST", "/api/friends/follow/"+alice.ID, bob.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("bob follow alice: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	if w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("alice follow bob: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	if w := authedRequest(t, router, "POST", "/api/posts/"+post1+"/react", bob.Token, `{"value":1}`); w.Code != http.StatusOK {
		t.Fatalf("bob like post1: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	if w := authedRequest(t, router, "POST", "/api/posts/"+post2+"/react", bob.Token, `{"value":1}`); w.Code != http.StatusOK {
		t.Fatalf("bob like post2: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w := authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/gamification", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("alice gamification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	stats := decodeGamification(t, w.Body.Bytes())
	wantAlice := gamificationResponse{
		Level:     2,
		Total:     7,
		Posts:     metricResponse{Level: 1, Count: 3},
		Likes:     metricResponse{Level: 1, Count: 2},
		Followers: metricResponse{Level: 0, Count: 1},
		Following: metricResponse{Level: 0, Count: 1},
	}
	if stats != wantAlice {
		t.Fatalf("alice stats wrong:\n got  %+v\n want %+v", stats, wantAlice)
	}

	w = authedRequest(t, router, "GET", "/api/users/"+bob.ID+"/gamification", bob.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("bob gamification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	stats = decodeGamification(t, w.Body.Bytes())
	wantBob := gamificationResponse{
		Level:     1,
		Total:     2,
		Posts:     metricResponse{Level: 0, Count: 0},
		Likes:     metricResponse{Level: 0, Count: 0},
		Followers: metricResponse{Level: 0, Count: 1},
		Following: metricResponse{Level: 0, Count: 1},
	}
	if stats != wantBob {
		t.Fatalf("bob stats wrong:\n got  %+v\n want %+v", stats, wantBob)
	}
}

func TestGamification_CountsReplyLikes(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "grl-a", "grl-a@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "grl-b", "grl-b@test.com", "StrongPass123!")

	postID := createPost(t, router, alice.Token, "a post")
	replyID := createComment(t, router, alice.Token, postID, "a reply by alice")

	if w := authedRequest(t, router, "POST",
		"/api/posts/"+postID+"/comments/"+replyID+"/react", bob.Token, `{"value":1}`); w.Code != http.StatusOK {
		t.Fatalf("bob like alice reply: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w := authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/gamification", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("gamification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	stats := decodeGamification(t, w.Body.Bytes())
	if stats.Likes.Count != 1 || stats.Posts.Count != 1 || stats.Total != 2 {
		t.Fatalf("expected posts=1 likes=1 total=2, got %+v", stats)
	}
}

type leaderboardEntry struct {
	ID       string               `json:"id"`
	Username string               `json:"username"`
	Avatar   string               `json:"avatar"`
	Stats    gamificationResponse `json:"stats"`
}

func TestLeaderboard_ListsUsersWithStats(t *testing.T) {
	router, _ := SetupTestEnv()
	alice := registerAndLogin(t, router, "boardalice", "boardalice@test.com", "StrongPass123!")
	bob := registerAndLogin(t, router, "boardbob", "boardbob@test.com", "StrongPass123!")

	avatarBody := `{"displayname":"Alice","username":"boardalice","bio":"","avatar":"https://example.com/alice.png"}`
	if w := authedRequest(t, router, "PUT", "/api/users/"+alice.ID, alice.Token, avatarBody); w.Code != http.StatusOK {
		t.Fatalf("set alice avatar: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	post := createPost(t, router, alice.Token, "hello")
	if w := authedRequest(t, router, "POST", "/api/posts/"+post+"/react", bob.Token, `{"value":1}`); w.Code != http.StatusOK {
		t.Fatalf("bob like alice post: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w := authedRequest(t, router, "GET", "/api/leaderboard", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("leaderboard: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	var entries []leaderboardEntry
	if err := json.Unmarshal(w.Body.Bytes(), &entries); err != nil {
		t.Fatalf("decode leaderboard: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 leaderboard entries, got %d - body: %s", len(entries), w.Body.String())
	}

	byID := map[string]leaderboardEntry{}
	for _, e := range entries {
		byID[e.ID] = e
	}

	gotAlice, ok := byID[alice.ID]
	if !ok {
		t.Fatal("alice missing from leaderboard")
	}
	if gotAlice.Username != "boardalice" || gotAlice.Avatar != "https://example.com/alice.png" {
		t.Errorf("alice entry = %+v, want username=boardalice avatar set", gotAlice)
	}
	if gotAlice.Stats.Total != 2 || gotAlice.Stats.Level != 1 {
		t.Errorf("alice stats = %+v, want total=2 level=1", gotAlice.Stats)
	}

	gotBob, ok := byID[bob.ID]
	if !ok {
		t.Fatal("bob missing from leaderboard")
	}
	if gotBob.Avatar != models.DefaultAvatar {
		t.Errorf("bob avatar = %q, want default avatar", gotBob.Avatar)
	}
	if gotBob.Stats.Total != 0 {
		t.Errorf("bob total = %d, want 0", gotBob.Stats.Total)
	}
}

func TestLeaderboard_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/leaderboard", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d - body: %s", w.Code, w.Body.String())
	}
}

func TestGamification_RequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "gameranon", "gameranon@test.com", "StrongPass123!")

	w := authedRequest(t, router, "GET", "/api/users/"+u.ID+"/gamification", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d - body: %s", w.Code, w.Body.String())
	}
}
