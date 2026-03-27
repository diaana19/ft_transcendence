package handlers

import (
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	// Routes
	api := router.Group("/api")
	{
		// Test route
		api.GET("/tweets", getTweets)
		api.POST("/tweets", createTweet)
	}
}

func getTweets(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "TODO: Implement getTweets",
	})
}

func createTweet(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "TODO: Implement createTweet",
	})
}
