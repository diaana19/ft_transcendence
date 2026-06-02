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

// Search godoc
// @Summary  Search across resources
// @Tags     search
// @Security BearerAuth
// @Produce  json
// @Param    q     query string true  "search query"
// @Param    type  query string false "type of resource to search" default(user)
// @Param    sort  query string false "field to sort by" default(username)
// @Param    order query string false "sort order (asc, desc)" default(asc)
// @Param    page  query int    false "page number" default(1)
// @Param    limit query int    false "max number of results per page (1-100)" default(20)
// @Success  200 {object} map[string]interface{}
// @Failure  400 {object} map[string]string
// @Failure  403 {object} map[string]string
// @Failure  500 {object} map[string]string
// @Router   /search [get]
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
