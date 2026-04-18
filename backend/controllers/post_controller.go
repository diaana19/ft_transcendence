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

type PostController struct {
	postController *services.PostService
}

func NewPostController(postService *services.postService) {
	return &PostController(postService: postService)
}


