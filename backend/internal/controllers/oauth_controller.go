package controllers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

type OAuthController struct {
	service     *services.OAuthService
	frontendURL string
}

func NewOAuthController(service *services.OAuthService, cfg *config.Config) *OAuthController {
	return &OAuthController{
		service:     service,
		frontendURL: cfg.FrontendURL,
	}
}

func (oc *OAuthController) OAuthLogin(c *gin.Context) {
	ctx := c.Request.Context()

	// Fail early if GITHUB_CLIENT_ID/SECRET aren't set. Otherwise we'd hand
	// the browser a github.com authorize URL with client_id= empty, which
	// just renders a GitHub 404 and confuses the user.
	if !oc.service.IsConfigured() {
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_not_configured")
		return
	}

	state, err := oc.service.GenerateState(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "could not start oauth flow",
		})
		return
	}

	url := oc.service.BuildAuthURL(state)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (oc *OAuthController) OAuthCallback(c *gin.Context) {
	ctx := c.Request.Context()

	code := c.Query("code")
	state := c.Query("state")

	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_denied")
		return
	}

	valid, err := oc.service.VerifyAndConsumeState(ctx, state)
	if err != nil {
		log.Printf("OAuth: state verification failed: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "state verification failed",
		})
		return
	}
	if !valid {
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=invalid_state")
		return
	}

	token, err := oc.service.ExchangeCodeForToken(ctx, code)
	if err != nil {
		log.Printf("OAuth: token exchange failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_failed")
		return
	}

	ghUser, err := oc.service.FetchGitHubUser(ctx, token)
	if err != nil {
		log.Printf("OAuth: fetch user failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_failed")
		return
	}

	user, err := oc.service.FindOrCreateUser(ctx, ghUser)
	if err != nil {
		log.Printf("OAuth: find or create user failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_failed")
		return
	}

	jwt, err := utils.GenerateJWT(user.ID, user.Username)
	if err != nil {
		log.Printf("OAuth: JWT generation failed: %v", err)
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/login?error=oauth_failed")
		return
	}

	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(
		"auth_token",
		jwt,
		24*3600,
		"/",
		"",
		false,
		true,
	)
	c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/")
}
