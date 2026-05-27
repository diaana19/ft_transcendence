package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/services"
)

type SearchController struct {
	service *services.SearchService
}

func NewSearchController(service *services.SearchService) *SearchController {
	return &SearchController{service: service}
}

func (sc *SearchController) Search(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusForbidden, gin.H{"error": "unauthorized"})
		return
	}

	q := c.Query("q")
	searchType := c.DefaultQuery("type", "user")
	sort := c.DefaultQuery("sort", "username")
	order := c.DefaultQuery("order", "asc")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	if q == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "query 'q' is required"})
		return
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 20
	}

	result, err := sc.service.Search(c.Request.Context(), userID.(string), q, searchType, sort, order, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, result)
}
