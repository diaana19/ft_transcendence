package controllers

import (
	"net/http"
	"strconv"

	"github.com/Transcendence/repositories"
	"github.com/gin-gonic/gin"
)

type MsgController struct {
	repo repositories.MessageRepository
}

func NewMsgController(repo repositories.MessageRepository) *MsgController {
	return &MsgController{repo: repo}
}

func (mc *MsgController) GetRoomMsg(c *gin.Context) {
	roomID := c.Param("roomId")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	since := c.Query("since")

	message, err := mc.repo.GetByRoomID(roomID, since, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": message})
}

func (mc *MsgController) GetReplies(c *gin.Context) {
	parentId := c.Param("messageId")

	replies, err := mc.repo.GetReplies(parentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": replies})
}
