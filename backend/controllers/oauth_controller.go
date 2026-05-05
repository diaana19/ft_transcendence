package controllers

import (
	"log"
	"net/http"

	"github.com/Transcendence/config"
	"github.com/Transcendence/services"
	"github.com/Transcendence/utils"
	"github.com/gin-gonic/gin"
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
		c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"login?error=oauth_denied")
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
