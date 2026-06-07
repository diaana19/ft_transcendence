package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
)

type FriendController struct {
	Service             *services.FriendService
	NotificationService *services.NotificationService
}

// SendFriendRequest godoc
// @Summary   Send a friend request to a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Target user ID to send the friend request to"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/request/{id} [post]
func (fc *FriendController) SendFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	targetID := c.Param("id")
	if err := fc.Service.SendRequest(userID.(string), targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[FriendRequest] sender userID=%s username=%q -> target userID=%s", userID, userUsername, targetID)
	_ = fc.NotificationService.SendNotification(
		targetID,
		"",
		userID.(string),
		userUsername.(string),
		"friend_request",
		userUsername.(string)+" sent you a friend request",
		"",
	)
	c.JSON(http.StatusOK, gin.H{"message": "request sent"})
}

// AcceptFriend godoc
// @Summary   Accept a pending friend request
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Requester user ID whose friend request is being accepted"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/accept/{id} [post]
func (fc *FriendController) AcceptFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userUsername, exists := c.Get("username")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requesterID := c.Param("id")
	if err := fc.Service.AcceptRequest(userID.(string), requesterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_ = fc.NotificationService.SendNotification(
		requesterID,
		"",
		userID.(string),
		userUsername.(string),
		"friend_accept",
		userUsername.(string)+" accepted your friend request",
		"",
	)
	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

// FollowUser godoc
// @Summary   Follow a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Target user ID to follow"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/follow/{id} [post]
func (fc *FriendController) FollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	targetID := c.Param("id")
	if err := fc.Service.Follow(userID.(string), targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user followed"})
}

// UnfollowUser godoc
// @Summary   Unfollow a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Target user ID to unfollow"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/follow/{id} [delete]
func (fc *FriendController) UnfollowUser(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	targetID := c.Param("id")
	if err := fc.Service.Unfollow(userID.(string), targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user unfollowed"})
}

// RemoveFriend godoc
// @Summary   Remove an existing friend
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Friend user ID to remove"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/{id} [delete]
func (fc *FriendController) RemoveFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	targetID := c.Param("id")
	if err := fc.Service.RemoveFriend(userID.(string), targetID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend removed"})
}

// RejectFriendRequest godoc
// @Summary   Reject a pending friend request
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Requester user ID whose friend request is being rejected"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /friends/reject/{id} [post]
func (fc *FriendController) RejectFriendRequest(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requesterID := c.Param("id")
	if err := fc.Service.RejectRequest(userID.(string), requesterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "request rejected"})
}

// GetFollowers godoc
// @Summary   List the followers of a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "User ID whose followers are listed"
// @Success   200 {array}  models.UserResponse
// @Failure   500 {object} map[string]string
// @Router    /users/{id}/followers [get]
func (fc *FriendController) GetFollowers(c *gin.Context) {
	userID := c.Param("id")
	followers, err := fc.Service.GetFollowers(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, len(followers))
	for i, u := range followers {
		responses[i] = u.ToResponse()
	}
	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetFollowing godoc
// @Summary   List the users a user is following
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "User ID whose following list is returned"
// @Success   200 {array}  models.UserResponse
// @Failure   500 {object} map[string]string
// @Router    /users/{id}/following [get]
func (fc *FriendController) GetFollowing(c *gin.Context) {
	userID := c.Param("id")
	following, err := fc.Service.GetFollowing(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, len(following))
	for i, u := range following {
		responses[i] = u.ToResponse()
	}
	c.JSON(http.StatusOK, gin.H{"data": responses})
}

// GetFriends godoc
// @Summary   List the friends of a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "User ID whose friends are listed"
// @Success   200 {array}  models.UserResponse
// @Failure   500 {object} map[string]string
// @Router    /users/{id}/friends [get]
func (fc *FriendController) GetFriends(c *gin.Context) {
	userID := c.Param("id")
	friends, err := fc.Service.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]models.UserResponse, len(friends))
	for i, u := range friends {
		responses[i] = u.ToResponse()
	}
	c.JSON(http.StatusOK, gin.H{"data": responses})
}
