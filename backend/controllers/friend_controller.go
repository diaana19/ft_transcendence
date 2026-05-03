package controllers

import (
	"log"
	"net/http"

	"github.com/Transcendence/services"

	"github.com/gin-gonic/gin"
)

type FriendController struct {
	Service             *services.FriendService
	NotificationService *services.NotificationService
}

func (fc *FriendController) SendFriendRequest(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	userUsername := c.MustGet("username").(string)
	targetID := c.Param("id")

	err := fc.Service.SendRequest(userID, targetID)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[FriendRequest] sender userID=%s username=%q -> target userID=%s", userID, userUsername, targetID)
	fc.NotificationService.SendNotification(
		targetID,
		userUsername,
		userID,
		userUsername,
		"friend_request",
		userUsername+" sent you a friend request",
	)
	c.JSON(200, gin.H{"message": "request sent"})
}

func (fc *FriendController) AcceptFriend(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	requesterID := c.Param("id")

	err := fc.Service.AcceptRequest(userID, requesterID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (fc *FriendController) FollowUser(c *gin.Context) {
	userID := c.MustGet("userID").(string)
	targetID := c.Param("id")

	err := fc.Service.Follow(userID, targetID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user followed"})
}
