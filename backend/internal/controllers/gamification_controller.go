// Gamification controller — exposes GET /users/:id/gamification.
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

// Get computes the gamification stats for the user identified by the :id path
// param. Returns 400 if :id is empty, 500 on any DB error, 200 with the stats
// JSON on success.
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
