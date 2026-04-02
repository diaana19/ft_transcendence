package main

import (
	"log"

	"github.com/Lord-Lucius/Transcendence/config"
	"github.com/Lord-Lucius/Transcendence/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	var DB, config, err = config.Load()
	if err == nil {
		log.Fatal(err)
	}

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

	DB.Ping()

	routes.SetupRoutes(router)

	if (config.ApiPort == "") {
		config.ApiPort = "8000"
	}

	router.Run(":" + config.ApiPort)
}