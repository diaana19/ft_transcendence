package main

import (
	"os"

	"github.com/Lord-Lucius/Transcendence/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Backend API is running",
			"status": "success",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	routes.SetupRoutes(router)

	port := os.Getenv("API_PORT")
	if (port == "") {
		port = "8000"
	}

	router.Run(":" + port)
}