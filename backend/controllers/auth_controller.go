package controllers

import (
	"log"
	"time"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService *services.AuthService
}

func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{authService: authService}
}

type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
}

func (ac *AuthController) RegisterUser(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	log.Printf("DEBUG: Received input: %+v\n", input)
	log.Printf("DEBUG: Password length: %d\n", len(input.Password))
	log.Printf("DEBUG: Password: %s\n", input.Password)

	parsedDate, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	if !utils.CheckUserAge(parsedDate) {
		c.JSON(400, gin.H{"error": "user must be older than 13"})
		return
	}

	if ok, errCode := utils.CheckPasswordFormat(input.Password, input.Username); !ok {
		passwordMessages := []string{
			"Password contains the username",
			"Password too short or missing character types",
		}
		c.JSON(400, gin.H{"error": passwordMessages[errCode-1]})
		return
	}

	user := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    input.Password,
		DateOfBirth: parsedDate,
	}

	response, err := ac.authService.CreateAuthUserService(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "creation didn't worked"})
		return
	}

	c.JSON(200, response)
}
