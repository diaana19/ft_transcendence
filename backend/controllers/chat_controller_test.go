package controllers

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
)

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
		c.Set("userID", userID)
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

func TestSendMessage_Created(t *testing.T) {
	svc := &fakeChatService{sendResp: &models.MessageResponse{ID: "x", SenderID: "a", RecipientID: "b", Content: "hi"}}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":"b","content":"hi"}`)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestSendMessage_BindError(t *testing.T) {
	r := newChatRouter(&fakeChatService{}, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 on bind error, got %d", w.Code)
	}
}

func TestSendMessage_EmptyContentError(t *testing.T) {
	svc := &fakeChatService{sendErr: services.ErrEmptyContent}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":"b","content":"hi"}`)
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for empty content, got %d", w.Code)
	}
}

func TestSendMessage_RecipientNotFoundError(t *testing.T) {
	svc := &fakeChatService{sendErr: services.ErrRecipientNotFound}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":"b","content":"hi"}`)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
}

func TestSendMessage_GenericError(t *testing.T) {
	svc := &fakeChatService{sendErr: errors.New("boom")}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "POST", "/chat/messages", `{"recipient_id":"b","content":"hi"}`)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500, got %d", w.Code)
	}
}

func TestPoll_Success(t *testing.T) {
	svc := &fakeChatService{pollResp: &models.PollResponse{Messages: []models.MessageResponse{}, NextCursor: ""}}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/poll", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPoll_PassesParsedLimit(t *testing.T) {
	svc := &fakeChatService{pollResp: &models.PollResponse{}}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/poll?since=abc&limit=10", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestPoll_InvalidLimitFallsBackToDefault(t *testing.T) {
	svc := &fakeChatService{pollResp: &models.PollResponse{}}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/poll?limit=notanumber", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
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

func TestListConversation_Success(t *testing.T) {
	svc := &fakeChatService{listResp: &models.PollResponse{}}
	r := newChatRouter(svc, "a")
	w := doJSON(r, "GET", "/chat/messages?with=b", "")
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
}

func TestListConversation_MissingWith(t *testing.T) {
	r := newChatRouter(&fakeChatService{}, "a")
	w := doJSON(r, "GET", "/chat/messages", "")
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "with") {
		t.Errorf("expected error to mention 'with', got %s", w.Body.String())
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

func TestParseLimit(t *testing.T) {
	cases := []struct {
		in   string
		want int
	}{
		{"", 0},
		{"abc", 0},
		{"42", 42},
	}
	for _, tc := range cases {
		if got := parseLimit(tc.in); got != tc.want {
			t.Errorf("parseLimit(%q)=%d, want %d", tc.in, got, tc.want)
		}
	}
}

