package controllers

import (
	"net/http"

	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
)

type NotificationController struct {
	notifService *services.NotificationService
}

func NewNotificationController(notifService *services.NotificationService) *NotificationController {
	return &NotificationController{notifService: notifService}
}

func (nc *NotificationController) GetUnread(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	notifs, err := nc.notifService.GetUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notifs, "total": len(notifs)})
}

func (nc *NotificationController) MarkAllRead(c *gin.Context) {
	userID := c.MustGet("user_id").(string)
	if err := nc.notifService.MarkAllRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all notification marked as read"})
}
