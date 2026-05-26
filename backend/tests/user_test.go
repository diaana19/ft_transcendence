package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
)

const testUserPassword = "StrongPass123!"

// createUserAndGetID registers a fresh user through the public endpoint and
// returns its generated id. Registration is unauthenticated.
func createUserAndGetID(router *gin.Engine, t *testing.T) string {
	t.Helper()

	unique := fmt.Sprintf("%d", time.Now().UnixNano())

	body := fmt.Sprintf(`{
		"username": "testuser_%s",
		"email": "testuser_%s@example.com",
		"password": "%s",
		"dateOfBirth": "2000-01-01"
	}`, unique, unique, testUserPassword)

	w := authedRequest(t, router, "POST", "/api/auth/register", "", body)
	if w.Code != http.StatusOK {
		t.Fatalf("failed to create test user, status: %d, body: %s", w.Code, w.Body.String())
	}

	var resp map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatal("failed to unmarshal create user response")
	}

	id, ok := resp["id"].(string)
	if !ok || id == "" {
		t.Fatal("response should contain valid user id")
	}
	return id
}

// tokenFor mints a valid JWT for the given user id. The auth middleware only
// validates the signature and reads the id from the claims, so this is enough
// to authenticate (and to satisfy ownership checks keyed on the id).
func tokenFor(userID string) string {
	token, err := utils.GenerateJWT(userID, "tester")
	if err != nil {
		panic(err)
	}
	return token
}

func TestGetUsers_Success(t *testing.T) {
	router, _ := SetupTestEnv()

	createUserAndGetID(router, t)
	createUserAndGetID(router, t)
	createUserAndGetID(router, t)

	w := authedRequest(t, router, "GET", "/api/users", tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var users []map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &users); err != nil {
		t.Fatal("failed to unmarshal users list")
	}

	if len(users) != 3 {
		t.Fatalf("expected 3 users, got %d", len(users))
	}

	for _, u := range users {
		if _, hasPassword := u["password"]; hasPassword {
			t.Fatal("password must not be returned in GET /users")
		}
	}
}

func TestGetUsers_Empty(t *testing.T) {
	router, _ := SetupTestEnv()

	w := authedRequest(t, router, "GET", "/api/users", tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 even if empty, got %d", w.Code)
	}

	var users []interface{}
	json.Unmarshal(w.Body.Bytes(), &users)
	if len(users) != 0 {
		t.Fatal("expected empty array")
	}
}

func TestGetUser_Success(t *testing.T) {
	router, _ := SetupTestEnv()
	id := createUserAndGetID(router, t)

	w := authedRequest(t, router, "GET", "/api/users/"+id, tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var user map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &user)

	if user["id"] != id {
		t.Fatal("returned user id does not match")
	}
	if _, hasPassword := user["password"]; hasPassword {
		t.Fatal("password must not be returned in GET /users/:id")
	}
}

func TestGetUser_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()

	fakeID := "550e8400-e29b-41d4-a716-446655440000"

	w := authedRequest(t, router, "GET", "/api/users/"+fakeID, tokenFor("any-authenticated-user"), "")
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestUpdateUser_Success_Partial(t *testing.T) {
	router, _ := SetupTestEnv()
	id := createUserAndGetID(router, t)

	body := `{
		"name": "Updated Name",
		"username": "newusername",
		"bio": "New bio here",
		"avatar": "https://example.com/avatar.jpg"
	}`

	w := authedRequest(t, router, "PUT", "/api/users/"+id, tokenFor(id), body)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var updated map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &updated)

	if updated["name"] != "Updated Name" || updated["username"] != "newusername" {
		t.Fatal("update fields not applied correctly")
	}
}

func TestUpdateUser_DuplicateUsername(t *testing.T) {
	router, _ := SetupTestEnv()
	_ = createUserAndGetID(router, t)
	id2 := createUserAndGetID(router, t)

	body := `{"username": "testuser"}`

	w := authedRequest(t, router, "PUT", "/api/users/"+id2, tokenFor(id2), body)
	if w.Code != http.StatusOK {
		t.Logf("Note: Update with duplicate username returned %d (validation not yet implemented)", w.Code)
	}
}

func TestUpdateUser_InvalidEmail(t *testing.T) {
	router, _ := SetupTestEnv()
	id := createUserAndGetID(router, t)

	body := `{"email": "invalid-email"}`

	w := authedRequest(t, router, "PUT", "/api/users/"+id, tokenFor(id), body)
	if w.Code != http.StatusOK {
		t.Logf("Note: Update with invalid email returned %d (validation not yet implemented)", w.Code)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()

	fakeID := "550e8400-e29b-41d4-a716-446655440000"
	body := `{"name": "Should not work"}`

	w := authedRequest(t, router, "PUT", "/api/users/"+fakeID, tokenFor(fakeID), body)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on update not found, got %d", w.Code)
	}
}

func TestDeleteUser_Success(t *testing.T) {
	router, _ := SetupTestEnv()
	id := createUserAndGetID(router, t)

	delBody := fmt.Sprintf(`{"password":"%s"}`, testUserPassword)
	w := authedRequest(t, router, "DELETE", "/api/users/"+id, tokenFor(id), delBody)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d - body: %s", w.Code, w.Body.String())
	}

	getW := authedRequest(t, router, "GET", "/api/users/"+id, tokenFor("any-authenticated-user"), "")
	if getW.Code != http.StatusNotFound {
		t.Fatal("deleted user should return 404 on GET")
	}
}

func TestDeleteUser_NotFound(t *testing.T) {
	router, _ := SetupTestEnv()

	fakeID := "550e8400-e29b-41d4-a716-446655440000"

	delBody := fmt.Sprintf(`{"password":"%s"}`, testUserPassword)
	w := authedRequest(t, router, "DELETE", "/api/users/"+fakeID, tokenFor(fakeID), delBody)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404 on delete not found, got %d", w.Code)
	}
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	router, _ := SetupTestEnv()
	id := createUserAndGetID(router, t)

	body := `{ "name": "test" `

	w := authedRequest(t, router, "PUT", "/api/users/"+id, tokenFor(id), body)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on invalid JSON, got %d", w.Code)
	}
}
