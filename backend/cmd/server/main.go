package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "ft_transcendence/backend/docs" // generated OpenAPI spec (run `make swagger`)
	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/redis"
	"ft_transcendence/backend/internal/routes"
)

// maxBodySize is the maximum request body size for JSON endpoints (1MB).
// File uploads have their own 25MB limit in the upload service.
const maxBodySize = 1 << 20 // 1MB

// BodySizeLimit middleware rejects requests with a body larger than maxBodySize.
func BodySizeLimit() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, maxBodySize)
		c.Next()
	}
}

// SecurityHeaders middleware adds common security headers including CSP.
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Content-Security-Policy restricts where resources can be loaded from.
		// - default-src 'self': only allow resources from our domain
		// - script-src 'self': only run scripts from our domain (blocks inline JS)
		// - style-src 'self' 'unsafe-inline': styles from our domain + inline styles (React needs this)
		// - img-src 'self' data: blob: images from our domain, data URIs, and blob URLs (for previews)
		// - media-src 'self' blob: videos/audio from our domain and blob URLs
		// - connect-src 'self' wss: ws: API calls to our domain + WebSocket connections
		// - frame-ancestors 'none': prevents page from being embedded in iframes (clickjacking)
		c.Header("Content-Security-Policy", "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: blob:; media-src 'self' blob:; connect-src 'self' wss: ws:; frame-ancestors 'none'")

		// Prevents browsers from MIME-sniffing a response away from declared content-type.
		// Stops a file labeled as image.jpg from being executed as JavaScript.
		c.Header("X-Content-Type-Options", "nosniff")

		// Blocks the page from being loaded in <iframe> or <frame>, preventing clickjacking attacks.
		c.Header("X-Frame-Options", "DENY")

		// Enables the browser's built-in XSS filter and blocks the page if an attack is detected.
		c.Header("X-XSS-Protection", "1; mode=block")

		// Controls how much referrer info is sent with requests:
		// - Same origin: full URL
		// - Cross origin: only the origin (not the full path)
		// - Downgrade (HTTPS→HTTP): nothing
		// Prevents leaking sensitive URL paths to external sites.
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		c.Next()
	}
}

// @title                       ft_transcendence API
// @version                     1.0
// @description                 Backend API for ft_transcendence (Pong + social platform).
// @BasePath                    /api
// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
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
		AllowOrigins:     []string{"https://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept", "X-Requested-With"},
		ExposeHeaders:    []string{"X-Request-Id", "X-RateLimit-Remaining", "Location"},
		AllowCredentials: true,
	}))
	router.Use(SecurityHeaders())
	router.Use(BodySizeLimit())

	api := router.Group("/api")
	api.GET("/", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"message": "Backend API is running",
			"status":  "success",
		})
	})

	// Health godoc
	// @Summary  Liveness probe
	// @Tags     health
	// @Produce  json
	// @Success  200  {object}  map[string]string
	// @Router   /health [get]
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

	// Swagger UI, reachable through the nginx `/swagger/` proxy entry (the backend
	// has no public port of its own). Reachable at <proxy>/swagger/index.html.
	// Spec built by `make swagger`.
	router.GET("/swagger/*any", ginSwagger.WrapHandler(
		swaggerFiles.Handler,
		ginSwagger.URL("/swagger/doc.json"),
	))

	if err := router.Run(":8080"); err != nil {
		log.Fatal("Server failed to start: ", err)
	}
}
