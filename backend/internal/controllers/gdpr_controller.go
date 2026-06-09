package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

// GDPRController handles the GDPR export and delete endpoints.
type GDPRController struct {
	gdprService *services.GDPRService
	mailService *services.MailService
}

// NewGDPRController creates a new GDPRController.
func NewGDPRController(gdprService *services.GDPRService, mailService *services.MailService) *GDPRController {
	return &GDPRController{gdprService: gdprService, mailService: mailService}
}

// ExportUserData exports all the data of the logged in user and emails a confirmation.
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

	if data.User.Email != "" {
		go func(email, username string) {
			subject := "Your Synk data export"
			body := "Hi " + username + ",\n\nThis confirms that you requested an export of your personal data " +
				"from Synk. The export was generated and downloaded from your browser.\n\n" +
				"If you did not request this, please secure your account and change your password.\n"
			if err := gc.mailService.SendMail([]string{email}, subject, body); err != nil {
				log.Printf("gdpr: export confirmation email failed for %s: %v", email, err)
			}
		}(data.User.Email, data.User.Username)
	}

	c.JSON(http.StatusOK, data)
}

// DeleteUserData deletes all the data of the logged in user and emails a confirmation.
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

	email, username, _ := gc.gdprService.GetUserContact(userID.(string))

	if err := gc.gdprService.DeleteUserData(userID.(string)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete user data"})
		return
	}

	if email != "" {
		go func(email, username string) {
			subject := "Your Synk data has been deleted"
			body := "Hi " + username + ",\n\nThis confirms that your personal data has been permanently " +
				"deleted from Synk, as you requested.\n\n" +
				"If you did not request this, please contact us immediately.\n"
			if err := gc.mailService.SendMail([]string{email}, subject, body); err != nil {
				log.Printf("gdpr: deletion confirmation email failed for %s: %v", email, err)
			}
		}(email, username)
	}

	c.JSON(http.StatusOK, gin.H{"message": "all user data has been permanently deleted"})
}
