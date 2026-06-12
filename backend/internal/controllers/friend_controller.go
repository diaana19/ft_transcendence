package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
)

// FriendController handles friend, follow and friend request endpoints.
type FriendController struct {
	Service             *services.FriendService
	NotificationService *services.NotificationService
}

// SendFriendRequest sends a friend request to the target user and notifies him.
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

// AcceptFriend accepts a pending friend request and notifies the requester.
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

// FollowUser makes the current user follow the target user.
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
	username, _ := c.Get("username")
	_ = fc.NotificationService.SendNotification(
		targetID, "",
		userID.(string), username.(string),
		"follow", username.(string)+" started following you",
		"",
	)
	c.JSON(http.StatusOK, gin.H{"message": "user followed"})
}

// UnfollowUser makes the current user stop following the target user.
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
	username, _ := c.Get("username")
	_ = fc.NotificationService.SendNotification(
		targetID, "",
		userID.(string), username.(string),
		"unfollow", username.(string)+" stopped following you",
		"",
	)
	c.JSON(http.StatusOK, gin.H{"message": "user unfollowed"})
}

// RemoveFriend removes the target user from the current user's friends.
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
	username, _ := c.Get("username")
	_ = fc.NotificationService.SendNotification(
		targetID, "",
		userID.(string), username.(string),
		"unfriend", username.(string)+" removed you from friends",
		"",
	)
	c.JSON(http.StatusOK, gin.H{"message": "friend removed"})
}

// RejectFriendRequest rejects a pending friend request from the requester.
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

// GetFriendStatus returns the relationship between the logged in user and the target.
// @Summary   Get the friend relationship status with a user
// @Tags      friends
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "Target user ID"
// @Success   200 {object} map[string]string "status: friends | pending_sent | pending_received | none"
// @Failure   401 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /friends/status/{id} [get]
func (fc *FriendController) GetFriendStatus(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	status, err := fc.Service.GetRelationship(userID.(string), c.Param("id"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": status})
}

// GetFollowers returns the list of users that follow the given user.
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

// GetFollowing returns the list of users that the given user follows.
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

// GetFriends returns the list of friends of the given user.
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
