package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

// NotificationController handles the notification endpoints.
type NotificationController struct {
	notifService *services.NotificationService
}

// NewNotificationController creates a new NotificationController.
func NewNotificationController(notifService *services.NotificationService) *NotificationController {
	return &NotificationController{notifService: notifService}
}

// GetUnread returns the unread notifications of the logged in user.
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

// MarkAllRead marks all notifications of the logged in user as read.
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

// MarkRead marks one notification as read by its id. Marking a notification that
// no longer exists is a no-op that reports success.
// @Summary   Mark a single notification as read
// @Tags      notifications
// @Security  BearerAuth
// @Produce   json
// @Param     id   path      string  true  "Notification ID"
// @Success   200  {object}  map[string]string
// @Failure   401  {object}  map[string]string
// @Failure   500  {object}  map[string]string
// @Router    /notification/{id}/read [patch]
func (nc *NotificationController) MarkRead(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	if err := nc.notifService.MarkRead(userID, c.Param("id")); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "notification marked as read"})
}
