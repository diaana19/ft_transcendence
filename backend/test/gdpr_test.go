package test

import (
	"encoding/json"
	"net/http"
	"testing"
)

func TestGDPR_ExportUserData(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "gdprexp", "gdprexp@test.com", "StrongPass123!")

	createPost(t, router, u.Token, "my first post")
	createPost(t, router, u.Token, "my second post")

	w := authedRequest(t, router, "GET", "/api/gdpr/export", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("export: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	var export struct {
		User struct {
			ID    string `json:"id"`
			Email string `json:"email"`
		} `json:"user"`
		Posts []struct {
			Content string `json:"content"`
		} `json:"posts"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &export); err != nil {
		t.Fatalf("decode export: %v", err)
	}
	if export.User.ID != u.ID {
		t.Fatalf("expected exported user id %s, got %s", u.ID, export.User.ID)
	}
	if len(export.Posts) != 2 {
		t.Fatalf("expected 2 exported posts, got %d", len(export.Posts))
	}
	if body := w.Body.String(); contains(body, "StrongPass123!") {
		t.Fatal("export must never include the plaintext password")
	}
}

func TestGDPR_ExportRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/gdpr/export", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func TestGDPR_DeleteUserData(t *testing.T) {
	router, _ := SetupTestEnv()
	u := registerAndLogin(t, router, "gdprdel", "gdprdel@test.com", "StrongPass123!")
	createPost(t, router, u.Token, "doomed post")

	w := authedRequest(t, router, "DELETE", "/api/gdpr/delete", u.Token, "")
	if w.Code != http.StatusOK {
		t.Fatalf("delete: expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	w = authedRequest(t, router, "GET", "/api/users/"+u.ID, tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("deleted user should return 404, got %d", w.Code)
	}

	w = authedRequest(t, router, "GET", "/api/posts/user/"+u.ID, tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusOK {
		t.Fatalf("posts by deleted user: expected 200, got %d", w.Code)
	}
	var resp struct {
		Data []any `json:"data"`
	}
	json.Unmarshal(w.Body.Bytes(), &resp)
	if len(resp.Data) != 0 {
		t.Fatalf("expected deleted user's posts to be gone, got %d", len(resp.Data))
	}
}

func TestGDPR_DeleteRequiresAuth(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "DELETE", "/api/gdpr/delete", "", "")
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401 without token, got %d", w.Code)
	}
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
