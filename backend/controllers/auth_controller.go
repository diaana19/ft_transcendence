package controllers

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Transcendence/models"
	"github.com/Transcendence/services"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
)

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

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

func (ac *AuthController) LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	user, err := ac.authService.LoginAuthUserService(input.Email, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}

func (ac *AuthController) RefreshToken(c *gin.Context) {
	tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	newToken, err := utils.RefreshToken(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": newToken})
}
