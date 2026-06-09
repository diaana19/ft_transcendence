package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

// GamificationController handles gamification stats and leaderboard endpoints.
type GamificationController struct {
	service *services.GamificationService
}

// NewGamificationController creates a new GamificationController.
func NewGamificationController(svc *services.GamificationService) *GamificationController {
	return &GamificationController{service: svc}
}

// Get returns the gamification stats of the given user.
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

// GetLeaderboard returns the global gamification leaderboard.
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
