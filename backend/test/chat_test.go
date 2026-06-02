package test

// This file holds the shared HTTP test helpers (registerAndLogin, authedRequest)
// used across the integration suite. Direct-message sending/history now runs over
// the chat WebSocket; see ws_test.go for that coverage.

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

type chatTestUser struct {
	ID    string
	Token string
}

func registerAndLogin(t *testing.T, router *gin.Engine, username, email, password string) chatTestUser { //nolint:unparam
	t.Helper()

	regBody := fmt.Sprintf(`{
		"username": "%s",
		"email": "%s",
		"password": "%s",
		"dateOfBirth": "2000-01-01"
	}`, username, email, password)

	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBufferString(regBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("register %s: expected 200, got %d - body: %s", username, w.Code, w.Body.String())
	}
	var regResp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &regResp); err != nil {
		t.Fatalf("decode register response: %v", err)
	}

	loginBody := fmt.Sprintf(`{"email":"%s","password":"%s"}`, email, password)
	req, _ = http.NewRequest("POST", "/api/auth/login", bytes.NewBufferString(loginBody))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("login %s: expected 200, got %d - body: %s", username, w.Code, w.Body.String())
	}
	var loginResp struct {
		Token string `json:"token"`
		User  struct {
			ID string `json:"id"`
		} `json:"user"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &loginResp); err != nil {
		t.Fatalf("decode login response: %v", err)
	}

	return chatTestUser{ID: loginResp.User.ID, Token: loginResp.Token}
}

func authedRequest(t *testing.T, router *gin.Engine, method, path, token, body string) *httptest.ResponseRecorder {
	t.Helper()
	var buf *bytes.Buffer
	if body == "" {
		buf = &bytes.Buffer{}
	} else {
		buf = bytes.NewBufferString(body)
	}
	req, _ := http.NewRequest(method, path, buf)
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	return w
}

// seedMessage inserts a direct message straight into the database. DM sending
// now happens over the chat WebSocket, so HTTP integration tests that just need
// a message to exist (e.g. search) write one directly instead.
func seedMessage(t *testing.T, db *gorm.DB, senderID, recipientID, content string) {
	t.Helper()
	id, err := uuid.NewV7()
	if err != nil {
		t.Fatalf("seed message uuid: %v", err)
	}
	msg := models.Message{
		ID:          id.String(),
		CreatedAt:   time.Now().UTC(),
		SenderID:    senderID,
		RecipientID: recipientID,
		Content:     content,
		Type:        "dm",
	}
	if err := db.Create(&msg).Error; err != nil {
		t.Fatalf("seed message: %v", err)
	}
}
