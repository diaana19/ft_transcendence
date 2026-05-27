/*
** File: gdpr_controller.go
** Description: Handles GDPR compliance HTTP requests
** Responsibilities:
** - Handle user data export requests
** - Handle user data deletion requests
** - Extract user identity from JWT token
 */

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
