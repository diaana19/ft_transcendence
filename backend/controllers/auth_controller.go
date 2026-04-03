package controllers

import (
	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// auth route
func RegisterUser(c *gin.Context, DB *gorm.DB) {
	var user models.User 
	err := c.BindJSON(&user);
	password_error_message := []string{"Password too short", "Password contains the user name or name"}

	if err != nil {
		c.JSON(400, gin.H{
			"error": "couldn't bind user",
		})
		return
	}

	has_password_check_worked, password_check_err_code := utils.CheckPasswordFormat(user.Password, user.Name, user.Username)

	if !utils.CheckEmailFormat(user.Email) {
		c.JSON(400, gin.H{
			"error": "invalid email format",
		})
		return
	} else if !has_password_check_worked {
		msg := password_error_message[password_check_err_code]
		c.JSON(400, gin.H{
			"error": msg,
		})
		return
	}

	response_service, err := services.CreateAuthUserService(&user, DB)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "authentification service didn't work well",
		})
		return
	}

	c.JSON(200, response_service)
}