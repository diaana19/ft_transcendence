package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
)

type ChatServicer interface {
	Send(senderID string, input models.CreateMessageInput) (*models.MessageResponse, error)
	Poll(userID, since string, limit int) (*models.PollResponse, error)
	ListConversation(userID, peerID, since string, limit int) (*models.PollResponse, error)
}

type ChatController struct {
	chatService ChatServicer
}

func NewChatController(chatService ChatServicer) *ChatController {
	return &ChatController{chatService: chatService}
}

// SendMessage godoc
// @Summary     Send a direct message (REST)
// @Description REST fallback for sending a DM; the primary path is the chat WebSocket (/api/ws/chat).
// @Tags        chat
// @Security    BearerAuth
// @Accept      json
// @Produce     json
// @Param       body  body      models.CreateMessageInput  true  "recipient id and content"
// @Success     201   {object}  models.MessageResponse
// @Failure     400   {object}  map[string]string  "empty content or invalid body"
// @Failure     404   {object}  map[string]string  "recipient not found"
// @Failure     500   {object}  map[string]string
// @Router      /chat/messages [post]
func (cc *ChatController) SendMessage(c *gin.Context) {
	userID := c.MustGet("user_id").(string)

	var input models.CreateMessageInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msg, err := cc.chatService.Send(userID, input)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrEmptyContent):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, services.ErrRecipientNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// Poll godoc
// @Summary     Poll for new messages
// @Description Long-poll style fetch of messages newer than a cursor across all the caller's conversations.
// @Tags        chat
// @Security    BearerAuth
// @Produce     json
// @Param       since  query     string  false  "cursor: return messages after this message id"
// @Param       limit  query     int     false  "page size (default 50, max 200)"
// @Success     200    {object}  models.PollResponse
// @Failure     500    {object}  map[string]string
// @Router      /chat/poll [get]
func (cc *ChatController) Poll(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	since := c.Query("since")
	limit := parseLimit(c.Query("limit"))

	resp, err := cc.chatService.Poll(userID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// ListConversation godoc
// @Summary     List a direct-message conversation
// @Description Returns the message history between the caller and a peer.
// @Tags        chat
// @Security    BearerAuth
// @Produce     json
// @Param       with   query     string  true   "peer user id"
// @Param       since  query     string  false  "cursor: return messages after this message id"
// @Param       limit  query     int     false  "page size (default 50, max 200)"
// @Success     200    {object}  models.PollResponse
// @Failure     400    {object}  map[string]string  "missing 'with' query parameter"
// @Failure     500    {object}  map[string]string
// @Router      /chat/messages [get]
func (cc *ChatController) ListConversation(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	peerID := c.Query("with")
	if peerID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing 'with' query parameter"})
		return
	}
	since := c.Query("since")
	limit := parseLimit(c.Query("limit"))

	resp, err := cc.chatService.ListConversation(userID, peerID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func parseLimit(raw string) int {
	if raw == "" {
		return 0
	}
	n, err := strconv.Atoi(raw)
	if err != nil {
		return 0
	}
	return n
}
