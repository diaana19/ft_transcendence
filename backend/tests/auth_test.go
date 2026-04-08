package tests

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegisterUser_Success(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{
		"username": "orion123",
		"email": "orion@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2005-06-15"
	}`

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	bodyStr := w.Body.String()
	if strings.Contains(bodyStr, "password") || strings.Contains(bodyStr, "StrongPass123!") {
		t.Fatal("password must never be returned in the response")
	}
}

func TestRegisterUser_InvalidJSON(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{ "username": "orion123", "email": "orion@test.com" `

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid JSON, got %d", w.Code)
	}
}

func TestRegisterUser_MissingRequiredFields(t *testing.T) {
	router, _ := SetupTestEnv()

	cases := []struct {
		name string
		body string
	}{
		{"missing username", `{"email":"test@test.com","password":"StrongPass123!","dateOfBirth":"2005-01-01"}`},
		{"missing email", `{"username":"testuser","password":"StrongPass123!","dateOfBirth":"2005-01-01"}`},
		{"missing password", `{"username":"testuser","email":"test@test.com","dateOfBirth":"2005-01-01"}`},
		{"missing dateOfBirth", `{"username":"testuser","email":"test@test.com","password":"StrongPass123!"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(tc.body)))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for %s, got %d", tc.name, w.Code)
			}
		})
	}
}

func TestRegisterUser_InvalidEmailFormat(t *testing.T) {
	router, _ := SetupTestEnv()

	invalidEmails := []string{"invalid-email", "test@", "@example.com", "test@example", "test@.com"}

	for _, email := range invalidEmails {
		body := fmt.Sprintf(`{
			"username": "testuser",
			"email": "%s",
			"password": "StrongPass123!",
			"dateOfBirth": "2005-01-01"
		}`, email)

		req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected 400 for invalid email '%s', got %d", email, w.Code)
		}
	}
}

func TestRegisterUser_Underage(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{
		"username": "younguser",
		"email": "young@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2014-04-08"
	}`

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for underage, got %d", w.Code)
	}
}

func TestRegisterUser_WeakPassword(t *testing.T) {
	router, _ := SetupTestEnv()

	testCases := []struct {
		name     string
		username string
		password string
	}{
		{"too short", "orion123", "Ab1!"},
		{"no uppercase", "orion123", "weakpass123!"},
		{"no lowercase", "orion123", "WEAKPASS123!"},
		{"no digit", "orion123", "WeakPass!!!"},
		{"contains username", "orion123", "OrionStrongPass123!"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body := fmt.Sprintf(`{
				"username": "%s",
				"email": "test@test.com",
				"password": "%s",
				"dateOfBirth": "2005-01-01"
			}`, tc.username, tc.password)

			req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("expected 400 for weak password case '%s', got %d - body: %s",
					tc.name, w.Code, w.Body.String())
			}
		})
	}
}

func TestRegisterUser_DuplicateEmail(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{
		"username": "user1",
		"email": "duplicate@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2000-01-01"
	}`

	req1, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	if w1.Code != http.StatusOK {
		t.Fatalf("first registration should succeed, got %d", w1.Code)
	}

	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for duplicate email, got %d - body: %s", w2.Code, w2.Body.String())
	}
}

func TestRegisterUser_DuplicateUsername(t *testing.T) {
	router, _ := SetupTestEnv()

	body1 := `{
		"username": "sameusername",
		"email": "user1@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2000-01-01"
	}`

	body2 := `{
		"username": "sameusername",
		"email": "user2@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2000-01-01"
	}`

	req1, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body1)))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body2)))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	if w2.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for duplicate username, got %d", w2.Code)
	}
}

func TestRegisterUser_ValidButEdgeCaseAge(t *testing.T) {
	router, _ := SetupTestEnv()

	body := `{
		"username": "thirteenuser",
		"email": "thirteen@test.com",
		"password": "StrongPass123!",
		"dateOfBirth": "2013-04-08"
	}`

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("user who is 13 years old should be allowed, got %d", w.Code)
	}
}
