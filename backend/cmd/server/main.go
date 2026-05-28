package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/redis"
	"ft_transcendence/backend/internal/routes"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	pdb, err := conf.Postgres.Connect()
	if err != nil {
		log.Fatal(err)
	}

	rdb, err := conf.Redis.Connect()
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	_ = router.SetTrustedProxies(nil)

	router.Use(cors.New(cors.Config{
		AllowOriginFunc:  func(_ string) bool { return true },
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With", "X-CSRF-Token"},
		ExposeHeaders:    []string{"X-Request-Id", "X-RateLimit-Remaining", "Location"},
		AllowCredentials: true,
	}))

	api := router.Group("/api")
	api.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Backend API is running",
			"status":  "success",
		})
	})

	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "OK",
		})
	})

	api.GET("/health/redis", func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 2*time.Second)
		defer cancel()
		if err := redis.RoundTrip(ctx, rdb, "health", "ping"); err != nil {
			c.JSON(503, gin.H{"status": "fail", "error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "OK"})
	})

	ctrl := routes.Wire(pdb, rdb, conf)
	routes.SetupRoutes(router, ctrl, rdb)

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
