package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

type GDPRController struct {
	gdprService *services.GDPRService
}

func NewGDPRController(gdprService *services.GDPRService) *GDPRController {
	return &GDPRController{gdprService: gdprService}
}

// ExportUserData godoc
// @Summary  Export the authenticated user's data
// @Tags     gdpr
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} map[string]interface{}
// @Failure  401 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router   /gdpr/export [get]
func (gc *GDPRController) ExportUserData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	data, err := gc.gdprService.ExportUserData(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not export user data"})
		return
	}

	c.JSON(http.StatusOK, data)
}

// DeleteUserData godoc
// @Summary  Permanently delete the authenticated user's data
// @Tags     gdpr
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router   /gdpr/delete [delete]
func (gc *GDPRController) DeleteUserData(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := gc.gdprService.DeleteUserData(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "all user data has been permanently deleted"})
}
