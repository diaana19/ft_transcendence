package main

import (
	"context"
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/Transcendence/config"
	"github.com/Transcendence/redis"
	"github.com/Transcendence/routes"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	DB, dberr := config.ConnectDB()
	if dberr != nil {
		log.Fatal(dberr)
	}

	router := gin.Default()
	router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(origin string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Cookie"},
		ExposeHeaders:    []string{"Content-Length", "Set-Cookie"},
		AllowCredentials: true,
	}))

	router.OPTIONS("/*path", func(c *gin.Context) {
		c.Status(204)
	})

	rdb, err := redis.InitRedis()
	if err != nil {
		log.Fatal(err)
	}
	if os.Getenv("DEBUG_REDIS") == "true" {
		ctx := context.Background()
		redis.Subscribe(ctx, rdb, "test-channel", func(message string) {
			log.Printf("Handler received: %s", message)
		})
		if err := redis.Publish(rdb, "test-channel", "Redis pub/sub is working"); err != nil {
			log.Printf("Publish error: %v\n", err)
		}
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

	if err := router.Run(":" + conf.ApiPort); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
