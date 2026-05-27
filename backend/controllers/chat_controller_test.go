package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	"github.com/Transcendence/models"
)

// fakeChatService lets the controller's error-mapping be tested in isolation.
// The success, bind-error, empty-content, recipient-not-found and missing-"with"
// cases are covered end-to-end through the /chat endpoints; only the generic
// service/DB error mappings (which cannot be triggered over HTTP) live here.
type fakeChatService struct {
	sendResp *models.MessageResponse
	sendErr  error
	pollResp *models.PollResponse
	pollErr  error
	listResp *models.PollResponse
	listErr  error
}

func (f *fakeChatService) Send(senderID string, input models.CreateMessageInput) (*models.MessageResponse, error) {
	return f.sendResp, f.sendErr
}

func (f *fakeChatService) Poll(userID, since string, limit int) (*models.PollResponse, error) {
	return f.pollResp, f.pollErr
}

func (f *fakeChatService) ListConversation(userID, peerID, since string, limit int) (*models.PollResponse, error) {
	return f.listResp, f.listErr
}

func newChatRouter(svc ChatServicer, userID string) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	cc := NewChatController(svc)
	r.Use(func(c *gin.Context) {
		c.Set("user_id", userID)
		c.Next()
	})
	r.POST("/chat/messages", cc.SendMessage)
	r.GET("/chat/messages", cc.ListConversation)
	r.GET("/chat/poll", cc.Poll)
	return r
}

func doJSON(r *gin.Engine, method, path, body string) *httptest.ResponseRecorder {
	var buf *bytes.Buffer
	if body == "" {
		buf = &bytes.Buffer{}
	} else {
		buf = bytes.NewBufferString(body)
	}
	req, _ := http.NewRequest(method, path, buf)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestSendMessage_GenericError(t *testing.T) {
	svc := &fakeChatService{sendErr: errors.New("boom")}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":"b","content":"hi"}`)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestPoll_ServiceError(t *testing.T) {
	svc := &fakeChatService{pollErr: errors.New("db")}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/poll", "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestListConversation_ServiceError(t *testing.T) {
	svc := &fakeChatService{listErr: errors.New("db")}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/messages?with=b", "")
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}
