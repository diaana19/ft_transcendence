package controllers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

type LoginInput struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password" binding:"required"`
}

type AuthController struct {
	authService  *services.AuthService
	twoFAService *services.TwoFAService
	mailService  *services.MailService
	rdb          *redis.Client
}

func NewAuthController(
	authService *services.AuthService,
	twoFAService *services.TwoFAService,
	mailService *services.MailService,
	rdb *redis.Client,
) *AuthController {
	return &AuthController{
		authService:  authService,
		twoFAService: twoFAService,
		mailService:  mailService,
		rdb:          rdb,
	}
}

type RegisterInput struct {
	Username    string `json:"username" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required"`
	DateOfBirth string `json:"dateOfBirth" binding:"required"`
}

// RegisterUser godoc
// @Summary   Register a new user
// @Description Create a local account after validating age and password strength
// @Tags      auth
// @Accept    json
// @Produce   json
// @Param     body body RegisterInput true "Registration details"
// @Success   200 {object} models.User
// @Failure   400 {object} map[string]string
// @Router    /auth/register [post]
func (ac *AuthController) RegisterUser(c *gin.Context) {
	var input RegisterInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	log.Printf("Register attempt: username=%s, ip=%s", input.Username, c.ClientIP())

	if !utils.CheckUsernameFormat(input.Username) {
		c.JSON(400, gin.H{"error": "username may only contain alphanumeric characters or single hyphens, " +
			"cannot begin or end with a hyphen, and must be 1-39 characters"})
		return
	}

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
			"Password must not contain the username",
			"Password too short (minimum 8 characters)",
			"Password must contain a lowercase letter",
			"Password must contain an uppercase letter",
			"Password must contain a digit",
			"Password must contain a special character",
		}
		c.JSON(400, gin.H{"error": passwordMessages[errCode]})
		return
	}

	password := input.Password
	user := models.User{
		Username:    input.Username,
		Email:       input.Email,
		Password:    &password,
		DateOfBirth: &parsedDate,
		Provider:    "local",
	}

	response, err := ac.authService.CreateAuthUserService(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if response.Email != "" {
		go func(email, username string) {
			subject := "Welcome to Synk"
			body := fmt.Sprintf(
				"Hi %s,\n\nYour Synk account has been created successfully. Welcome aboard!\n\n"+
					"If you didn't create this account, please contact us immediately.\n",
				username,
			)
			if err := ac.mailService.SendMail([]string{email}, subject, body); err != nil {
				log.Printf("welcome email failed for %s: %v", email, err)
			}
		}(response.Email, response.Username)
	}

	c.JSON(200, response)
}

// LoginUser godoc
// @Summary   Authenticate a user
// @Description Log in with email or username and password; may require a 2FA step
// @Tags      auth
// @Accept    json
// @Produce   json
// @Param     body body LoginInput true "Login credentials"
// @Success   200 {object} map[string]interface{}
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /auth/login [post]
func (ac *AuthController) LoginUser(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

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

	if user.TwoFAEnabled {
		pendingToken, tokenErr := ac.authService.CreatePendingLogin(user.ID, ac.rdb)
		if tokenErr != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "could not initiate 2FA flow",
			})
			return
		}

		log.Printf("🔐 2FA required for userID=%s, pending_token issued", user.ID)
		c.JSON(http.StatusOK, gin.H{
			"needs_2fa":     true,
			"pending_token": pendingToken,
		})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}
	c.SetCookie(
		"auth_token",
		token,
		int((24 * time.Hour).Seconds()),
		"/",
		"",
		true,
		true,
	)
	c.JSON(http.StatusOK, gin.H{"token": token, "user": user.ToResponse()})
}

// RefreshToken godoc
// @Summary   Refresh a JWT
// @Description Issue a new JWT from a valid bearer token
// @Tags      auth
// @Produce   json
// @Param     Authorization header string true "Bearer token"
// @Success   200 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Router    /auth/refresh [post]
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

type ForgotPasswordInput struct {
	Email string `json:"email" binding:"required,email"`
}

// ForgotPassword godoc
// @Summary   Request a password reset
// @Description Generate a new password for the account with the given email and send it to that address.
// @Description Always returns 200 so it can't be used to discover which emails are registered.
// @Tags      auth
// @Accept    json
// @Produce   json
// @Param     body body ForgotPasswordInput true "Account email"
// @Success   200 {object} map[string]string
// @Failure   400 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /auth/forgot-password [post]
func (ac *AuthController) ForgotPassword(c *gin.Context) {
	var input ForgotPasswordInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "a valid email is required"})
		return
	}

	const genericMsg = "if an account exists for that email, a new password has been sent"

	user, err := ac.authService.GetUserByEmail(input.Email)
	if err != nil {
		log.Printf("forgot-password: no account for email=%s ip=%s", input.Email, c.ClientIP())
		c.JSON(http.StatusOK, gin.H{"message": genericMsg})
		return
	}

	deliver := func(newPassword string) error {
		subject := "Your new ft_transcendence password"
		body := fmt.Sprintf(
			"Hi %s,\n\nYou requested a password reset. Your new password is:\n\n    %s\n\n"+
				"Please log in and change it from your profile as soon as possible.\n\n"+
				"If you didn't request this, change your password immediately.\n",
			user.Username, newPassword,
		)
		return ac.mailService.SendMail([]string{user.Email}, subject, body)
	}

	if err := ac.authService.ResetPassword(user.ID, deliver); err != nil {
		log.Printf("forgot-password: reset failed for userID=%s: %v", user.ID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not process password reset, please try again later"})
		return
	}

	log.Printf("forgot-password: new password sent to userID=%s email=%s", user.ID, user.Email)
	c.JSON(http.StatusOK, gin.H{"message": genericMsg})
}

// Me returns the authenticated user's profile. Used by the SPA on mount to
// resolve "who am I via cookie?" — the OAuth callback only sets the HttpOnly
// auth_token cookie, so JS can't read the JWT directly and needs an endpoint
// to recover the user identity from a cookie-only session.
// Me godoc
// @Summary   Get current user
// @Description Return the authenticated user's profile resolved from the session
// @Tags      auth
// @Security  BearerAuth
// @Produce   json
// @Success   200 {object} map[string]interface{}
// @Failure   401 {object} map[string]string
// @Failure   404 {object} map[string]string
// @Router    /auth/me [get]
func (ac *AuthController) Me(c *gin.Context) {
	userIDRaw, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	user, err := ac.authService.GetUserByID(userIDRaw.(string))
	if err != nil {
		c.SetCookie("auth_token", "", -1, "/", "", true, true)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": user.ToResponse()})
}

// LogoutUser godoc
// @Summary   Log out the current user
// @Description Blacklist the active token and clear the auth cookie
// @Tags      auth
// @Security  BearerAuth
// @Produce   json
// @Success   200 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /auth/logout [post]
func (ac *AuthController) LogoutUser(c *gin.Context) {
	tokenStr := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	if tokenStr == "" {
		if cookieToken, err := c.Cookie("auth_token"); err == nil {
			tokenStr = cookieToken
		}
	}

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
	c.SetCookie("auth_token", "", -1, "/", "", true, true)

	if err := ac.authService.LogoutAuthUserService(tokenStr, expiry, ac.rdb); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

type Verify2FAInput struct {
	PendingToken string `json:"pending_token" binding:"required"`
	Code         string `json:"code" binding:"required,len=6,numeric"`
}

// Verify2FA godoc
// @Summary   Verify a 2FA code
// @Description Complete a pending login by validating the user's TOTP code
// @Tags      auth
// @Accept    json
// @Produce   json
// @Param     body body Verify2FAInput true "Pending token and 2FA code"
// @Success   200 {object} map[string]interface{}
// @Failure   400 {object} map[string]string
// @Failure   401 {object} map[string]string
// @Failure   500 {object} map[string]string
// @Router    /auth/2fa/verify [post]
func (ac *AuthController) Verify2FA(c *gin.Context) {
	var input Verify2FAInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	userID, err := ac.authService.PeekPendingLogin(input.PendingToken, ac.rdb)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	valid, err := ac.twoFAService.ValidateCode(userID, input.Code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": msgInvalid2FACode})
		return
	}
	if !valid {
		log.Printf("2FA verify failed: userID=%s, ip=%s", userID, c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": msgInvalid2FACode})
		return
	}

	ac.authService.ConsumePendingLogin(input.PendingToken, ac.rdb)

	user, err := ac.authService.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not fetch user"})
		return
	}

	token, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie(
		"auth_token",
		token,
		int((24 * time.Hour).Seconds()),
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user.ToResponse(),
	})
}
