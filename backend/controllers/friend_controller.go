package controllers

import (
    "net/http"
    "strconv"

    "github.com/Transcendence/services"

    "github.com/gin-gonic/gin"
)

type FriendController struct {
    Service *services.FriendService
}

func (fc *FriendController) SendFriendRequest(c *gin.Context) {
    userID := c.GetUint("userID")
    targetID, _ := strconv.Atoi(c.Param("id"))

    err := fc.Service.SendRequest(userID, uint(targetID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "request sent"})
}

func (fc *FriendController) AcceptFriend(c *gin.Context) {
    userID := c.GetUint("userID")
    requesterID, _ := strconv.Atoi(c.Param("id"))

    err := fc.Service.AcceptRequest(userID, uint(requesterID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "friend request accepted"})
}

func (fc *FriendController) FollowUser(c *gin.Context) {
    userID := c.GetUint("userID")
    targetID, _ := strconv.Atoi(c.Param("id"))

    err := fc.Service.Follow(userID, uint(targetID))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "user followed"})
}
