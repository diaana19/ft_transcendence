package controllers

import (
	"log"

	"net/http"

	"github.com/Transcendence/services"
	"github.com/Transcendence/models"
	"github.com/gin-gonic/gin"
)

type FriendController struct {
	Service             *services.FriendService
	NotificationService *services.NotificationService
}

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
	fc.NotificationService.SendNotification(
		targetID,
		userUsername.(string),
		userID.(string),
		userUsername.(string),
		"friend_request",
		userUsername.(string)+" sent you a friend request",
	)
	c.JSON(http.StatusOK, gin.H{"message": "request sent"})
}

func (fc *FriendController) AcceptFriend(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	requesterID := c.Param("id")
	if err := fc.Service.AcceptRequest(userID.(string), requesterID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

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
