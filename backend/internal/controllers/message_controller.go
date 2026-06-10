package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
)

// MessageController handles the direct message listing endpoints.
type MessageController struct {
	msgRepo     repositories.MessageRepository
	userService *services.UserService
}

// NewMessageController creates a new MessageController.
func NewMessageController(
	msgRepo repositories.MessageRepository,
	userService *services.UserService,
) *MessageController {
	return &MessageController{msgRepo: msgRepo, userService: userService}
}

// GetConversations returns the users the logged user has a direct chat with.
// @Summary   List the chat partners of the logged user
// @Tags      messages
// @Security  BearerAuth
// @Produce   json
// @Success   200 {object} map[string]interface{}
// @Failure   401 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /conversations [get]
func (mc *MessageController) GetConversations(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	peerIDs, err := mc.msgRepo.ListConversationPartnerIDs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, 0, len(peerIDs))
	for _, peerID := range peerIDs {
		user, err := mc.userService.GetUser(peerID)
		if err != nil {
			continue
		}
		responses = append(responses, user.ToResponse())
	}

	c.JSON(http.StatusOK, gin.H{"data": responses})
}
