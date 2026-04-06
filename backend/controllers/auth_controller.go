package controllers

import (
	"log"
	"net/http"
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

// Input JSON struct séparé du modèle DB
type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"` // string pour parser
}

func (ac *AuthController) RegisterUser(c *gin.Context) {
	var input RegisterInput

	// Bind et validation automatique par Gin
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	log.Printf("DEBUG: Received input: %+v\n", input)
	log.Printf("DEBUG: Password length: %d\n", len(input.Password))
	log.Printf("DEBUG: Password: %s\n", input.Password)

	// Parse la date
	parsedDate, err := time.Parse("2006-01-02", input.DateOfBirth)
	if err != nil {
		c.JSON(400, gin.H{"error": "invalid date format (expected YYYY-MM-DD)"})
		return
	}

	// Vérifie âge
	if !utils.CheckUserAge(parsedDate) {
		c.JSON(400, gin.H{"error": "user must be older than 13"})
		return
	}

	// Vérifie password
	if ok, errCode := utils.CheckPasswordFormat(input.Password, input.Username); !ok {
		passwordMessages := []string{"Password too short", "Password contains the user name or name", "Password need to contain at least maj, min,"}
		c.JSON(400, gin.H{"error": passwordMessages[errCode-1]})
		return
	}

	// Crée le user à passer au service
	user := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    input.Password,
		DateOfBirth: parsedDate,
	}

	// Appelle le service pour créer le user
	response, err := ac.authService.CreateAuthUserService(&user)
	if err != nil {
		c.JSON(500, gin.H{"error": "creation didn't worked"})
		return
	}

	c.JSON(http.StatusCreated, response)
}
