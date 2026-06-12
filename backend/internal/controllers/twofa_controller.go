package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

// TwoFAController handles the 2FA setup, enable and disable endpoints.
type TwoFAController struct {
	service     *services.TwoFAService
	userService *services.UserService
	mailService *services.MailService
}

// NewTwoFAController creates a new TwoFAController.
func NewTwoFAController(
	service *services.TwoFAService,
	userService *services.UserService,
	mailService *services.MailService,
) *TwoFAController {
	return &TwoFAController{service: service, userService: userService, mailService: mailService}
}

// notify2FAChange emails the user that 2FA was enabled or disabled.
func (tc *TwoFAController) notify2FAChange(userID, action string) {
	user, err := tc.userService.GetUser(userID)
	if err != nil || user.Email == "" {
		return
	}
	go func(email, username string) {
		subject := "Two-factor authentication " + action
		body := fmt.Sprintf(
			"Hi %s,\n\nTwo-factor authentication was %s on your Synk account.\n\n"+
				"If you didn't do this, secure your account immediately and contact us.\n",
			username, action,
		)
		if err := tc.mailService.SendMail([]string{email}, subject, body); err != nil {
			log.Printf("2fa %s email failed for %s: %v", action, email, err)
		}
	}(user.Email, user.Username)
}

// Setup starts the 2FA setup and returns the provisioning data.
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

// Enable2FAInput holds the TOTP code used to enable 2FA.
type Enable2FAInput struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

// Enable turns on 2FA after checking the TOTP code.
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
	tc.notify2FAChange(userID, "enabled")
	c.JSON(http.StatusOK, gin.H{"message": "2FA enabled successfully"})
}

// Disable2FAInput holds the TOTP code used to disable 2FA.
type Disable2FAInput struct {
	Code string `json:"code" binding:"required,len=6,numeric"`
}

// Disable turns off 2FA after checking the TOTP code. When 2FA is already
// disabled it is a no-op that reports success, so a repeated click never errors.
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

	if user, err := tc.userService.GetUser(userID); err == nil && !user.TwoFAEnabled {
		c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
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

	tc.notify2FAChange(userID, "disabled")
	c.JSON(http.StatusOK, gin.H{"message": "2FA disabled"})
}
