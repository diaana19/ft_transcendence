package test

import (
	"encoding/json"
	"net/http"
	"testing"
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

	// Bob follows Alice (alice gains a follower).
	if w := authedRequest(t, router, "POST", "/api/friends/follow/"+alice.ID, bob.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("bob follow alice: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	// Alice follows Bob (alice gains a "following").
	if w := authedRequest(t, router, "POST", "/api/friends/follow/"+bob.ID, alice.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("alice follow bob: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	// Bob likes two of alice's posts (alice gains 2 received likes).
	if w := authedRequest(t, router, "POST", "/api/posts/"+post1+"/like", bob.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("bob like post1: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	if w := authedRequest(t, router, "POST", "/api/posts/"+post2+"/like", bob.Token, ""); w.Code != http.StatusOK {
		t.Fatalf("bob like post2: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	// Alice: 3 posts + 2 likes received + 1 follower + 1 following = 7
	// 2^2 = 4 <= 7 < 8 = 2^3, so level = 2.
	w := authedRequest(t, router, "GET", "/api/users/"+alice.ID+"/gamification", alice.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("alice gamification: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}
	stats := decodeGamification(t, w.Body.Bytes())
	wantAlice := gamificationResponse{
		Level:     2,
		Total:     7,
		Posts:     metricResponse{Level: 1, Count: 3}, // 2^1 <= 3 < 2^2
		Likes:     metricResponse{Level: 1, Count: 2}, // 2^1 = 2
		Followers: metricResponse{Level: 0, Count: 1},
		Following: metricResponse{Level: 0, Count: 1},
	}
	if stats != wantAlice {
		t.Fatalf("alice stats wrong:\n got  %+v\n want %+v", stats, wantAlice)
	}

	// Bob: 0 posts + 0 likes received + 1 follower (alice) + 1 following (alice) = 2 → level 1.
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

	// Give alice an avatar so the non-nil branch of stringVal is exercised; bob
	// keeps the default empty avatar (nil → "").
	avatarBody := `{"displayname":"Alice","username":"boardalice","bio":"","avatar":"https://example.com/alice.png"}`
	if w := authedRequest(t, router, "PUT", "/api/users/"+alice.ID, alice.Token, avatarBody); w.Code != http.StatusOK {
		t.Fatalf("set alice avatar: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	// Alice: 1 post + 1 like received = total 2 → level 1. Bob stays at 0.
	post := createPost(t, router, alice.Token, "hello")
	if w := authedRequest(t, router, "POST", "/api/posts/"+post+"/like", bob.Token, ""); w.Code != http.StatusOK {
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
	if gotBob.Avatar != "" {
		t.Errorf("bob avatar = %q, want empty string", gotBob.Avatar)
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

	// No Authorization header → handled by the protected group's AuthMiddleware.
	w := authedRequest(t, router, "GET", "/api/users/"+u.ID+"/gamification", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d - body: %s", w.Code, w.Body.String())
	}
}
