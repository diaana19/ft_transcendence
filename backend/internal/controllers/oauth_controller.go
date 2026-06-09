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

// @Summary   Start GitHub OAuth login
// @Description Redirect the user to GitHub's authorization page to begin the OAuth flow
// @Tags      oauth
// @Produce   json
// @Success   307 "Redirect to GitHub authorization URL"
// @Failure   500 {object} map[string]string
// @Router    /auth/oauth/github/login [get]
func (oc *OAuthController) OAuthLogin(c *gin.Context) {
	ctx := c.Request.Context()

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

// @Summary   Handle GitHub OAuth callback
// @Description Validate state, exchange the code for a token, resolve the user,
// @Description set the auth cookie, and redirect to the frontend
// @Tags      oauth
// @Produce   json
// @Param     code  query string true  "Authorization code returned by GitHub"
// @Param     state query string false "Anti-CSRF state token"
// @Success   307 "Redirect to the frontend application"
// @Failure   500 {object} map[string]string
// @Router    /auth/oauth/github/callback [get]
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
		true,
		true,
	)
	c.Redirect(http.StatusTemporaryRedirect, oc.frontendURL+"/")
}
