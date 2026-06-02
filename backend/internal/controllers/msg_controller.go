package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/repositories"
)

type MsgController struct {
	repo repositories.MessageRepository
}

func NewMsgController(repo repositories.MessageRepository) *MsgController {
	return &MsgController{repo: repo}
}

// GetRoomMsg godoc
// @Summary  List messages in a room
// @Tags     chat
// @Security BearerAuth
// @Produce  json
// @Param    roomId path  string true  "room ID"
// @Param    since  query string false "return messages after this cursor/timestamp"
// @Param    limit  query int    false "max number of messages to return" default(50)
// @Success  200 {object} map[string]interface{}
// @Failure  500 {object} map[string]string
// @Router   /rooms/{roomId}/messages [get]
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

// GetReplies godoc
// @Summary  List replies to a message
// @Tags     chat
// @Security BearerAuth
// @Produce  json
// @Param    messageId path string true "parent message ID"
// @Success  200 {object} map[string]interface{}
// @Failure  500 {object} map[string]string
// @Router   /messages/{messageId}/replies [get]
func (mc *MsgController) GetReplies(c *gin.Context) {
	parentId := c.Param("messageId")

	replies, err := mc.repo.GetReplies(parentId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": replies})
}
