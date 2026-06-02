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

// Setup godoc
// @Summary  Begin 2FA setup and return provisioning data
// @Tags     2fa
// @Security BearerAuth
// @Produce  json
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Router   /2fa/setup [post]
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

// Enable godoc
// @Summary  Enable 2FA using a TOTP code
// @Tags     2fa
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    input body controllers.Enable2FAInput true "TOTP code"
// @Success  200 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Router   /2fa/enable [post]
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

// Disable godoc
// @Summary  Disable 2FA using a TOTP code
// @Tags     2fa
// @Security BearerAuth
// @Accept   json
// @Produce  json
// @Param    input body controllers.Disable2FAInput true "TOTP code"
// @Success  200 {object} map[string]string
// @Failure  400 {object} map[string]string
// @Failure  401 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router   /2fa/disable [post]
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
