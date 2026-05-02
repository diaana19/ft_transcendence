package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
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

func (cc *ChatController) SendMessage(c *gin.Context) {
	userID := c.MustGet("userID").(string)

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

func (cc *ChatController) Poll(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	since := c.Query("since")
	limit := parseLimit(c.Query("limit"))

	resp, err := cc.chatService.Poll(userID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

func (cc *ChatController) ListConversation(c *gin.Context) {
	userID := c.MustGet("userID").(string)
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
