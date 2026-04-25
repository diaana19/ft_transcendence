package controllers

import (
	"net/http"

	"github.com/Transcendence/config"
	"github.com/Transcendence/services"
	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	service		*services.OAuthService
	frontendURL	string
}

func NewOAuthController(service *services.OAuthService, cfg *config.Config) *OAuthController {
	return &OAuthController{
		service: service,
		frontendURL: cfg.FrontendURL,
	}
}

func (oc *OAuthController) OAuthLogin(c * gin.Context) {
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
