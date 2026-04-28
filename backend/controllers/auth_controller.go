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
	"github.com/redis/go-redis/v9"
)

// Identifier can be either email or username
type LoginInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type AuthController struct {
	authService *services.AuthService
	rdb         *redis.Client
}

func NewAuthController(authService *services.AuthService, rdb *redis.Client) *AuthController {
	return &AuthController{authService: authService, rdb: rdb}
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
			"Password too short",
			"Password don't contains lowercase",
			"Passowrd don't contains uppercase",
			"Password don't contains digit",
			"Password don't contains specials",
		}
		c.JSON(400, gin.H{"error": passwordMessages[errCode]})
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

func (ac *AuthController) LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	// determine identifier
	identifier := input.Email
	if identifier == "" {
		identifier = input.Username
	}
	if identifier == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email or username required"})
		return
	}

	log.Printf("🔐 Login attempt: identifier=%s, ip=%s", identifier, c.ClientIP())

	user, err := ac.authService.LoginAuthUserService(identifier, input.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	log.Printf("✅ Login success: userID=%s, ip=%s, username=%s", user.ID, c.ClientIP(), user.Username)
	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token, "user": user})
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

// logout use redis by putting the token in redis db

func (ac *AuthController) LogoutUser(c *gin.Context) {
	tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")
	if tokenStr == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	claims, err := utils.ValidateJWT(tokenStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	expiry := time.Until(claims.ExpiresAt.Time)
	if err := ac.authService.LogoutAuthUserService(tokenStr, expiry, ac.rdb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}
