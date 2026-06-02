package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

type NotificationController struct {
	notifService *services.NotificationService
}

func NewNotificationController(notifService *services.NotificationService) *NotificationController {
	return &NotificationController{notifService: notifService}
}

// GetUnread godoc
// @Summary   List unread notifications
// @Tags      notifications
// @Security  BearerAuth
// @Produce   json
// @Success   200  {object}  map[string]interface{}  "data: notifications, total: count"
// @Failure   401  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /notification [get]
func (nc *NotificationController) GetUnread(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	notifs, err := nc.notifService.GetUnread(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": notifs, "total": len(notifs)})
}

// MarkAllRead godoc
// @Summary   Mark all notifications as read
// @Tags      notifications
// @Security  BearerAuth
// @Produce   json
// @Success   200  {object}  map[string]string
// @Failure   401  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /notification/read [patch]
func (nc *NotificationController) MarkAllRead(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	if err := nc.notifService.MarkAllRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "all notification marked as read"})
}
