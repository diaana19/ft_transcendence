package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

type GamificationController struct {
	service *services.GamificationService
}

func NewGamificationController(svc *services.GamificationService) *GamificationController {
	return &GamificationController{service: svc}
}

// @Summary   Get gamification stats for a user
// @Tags      gamification
// @Security  BearerAuth
// @Produce   json
// @Param     id path string true "user id"
// @Success   200 {object} map[string]interface{}
// @Failure   400 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /users/{id}/gamification [get]
func (gc *GamificationController) Get(c *gin.Context) {
	userID := c.Param("id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing user id"})
		return
	}
	stats, err := gc.service.Compute(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to compute gamification stats"})
		return
	}
	c.JSON(http.StatusOK, stats)
}

// @Summary   Global gamification leaderboard
// @Tags      gamification
// @Security  BearerAuth
// @Produce   json
// @Success   200 {array} services.LeaderboardEntry
// @Failure   500 {object} map[string]string
// @Router    /leaderboard [get]
func (gc *GamificationController) GetLeaderboard(c *gin.Context) {
	entries, err := gc.service.Leaderboard()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch leaderboard"})
		return
	}
	c.JSON(http.StatusOK, entries)
}
