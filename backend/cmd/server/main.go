package main

import (
	"context"
	"log"

	"github.com/Transcendence/config"
	"github.com/Transcendence/redis"
	"github.com/Transcendence/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	var conf, err = config.Load()
	if err != nil {
		log.Fatal(err)
	}
	var DB, dberr = config.ConnectDB()
	if dberr != nil {
		log.Fatal(err)
	}

	var router *gin.Engine = gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	rdb, err := redis.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	redis.Subscribe(context.Background(), rdb, "test-channel", func(message string) {
		log.Printf("Handler received: %s", message)
	})

	if err := redis.Publish(rdb, "test-channel", "Redis pub/sub is working"); err != nil {
		log.Printf("Publish error: %v\n", err)
	}

	router.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Backend API is running",
			"status":  "success",
		})
	})

	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	routes.SetupRoutes(router, DB, rdb, conf)

	if conf.ApiPort == "" {
		conf.ApiPort = "8000"
	}

	router.Run(":" + conf.ApiPort)

}
