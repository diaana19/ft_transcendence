package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

type TwoFAController struct {
	service *services.TwoFAService
}

func NewTwoFAController(service *services.TwoFAService) *TwoFAController {
	return &TwoFAController{service: service}
}

func (tc *TwoFAController) Setup(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	resp, err := tc.service.GenerateSecret(userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, resp)
}

type Enable2FAInput struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

func (tc *TwoFAController) Enable(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID := userIDRaw.(string)

	var input Enable2FAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid code format (expected 6 digits)",
		})
		return
	}

	if err := tc.service.EnableTwoFA(userID, input.Code); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

type Disable2FAInput struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

func (tc *TwoFAController) Disable(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	userID := userIDRaw.(string)

	var input Disable2FAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	valid, err := tc.service.ValidateCode(userID, input.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !valid {
		c.JSON(http.StatusBadRequest, gin.H{"error": msgInvalid2FACode})
		return
	}

	if err := tc.service.DisableTwoFA(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}
