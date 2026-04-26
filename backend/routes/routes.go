package routes

import (
	"github.com/Transcendence/config"
	"github.com/Transcendence/controllers"
	"github.com/Transcendence/middleware"
	"github.com/Transcendence/repositories"
	"github.com/Transcendence/services"
	"github.com/Transcendence/socket"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func create_post_routes(api *gin.RouterGroup, DB *gorm.DB, rdb *redis.Client) {
	postRepo := repositories.NewPostRepository(DB)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService)

	posts := api.Group("/posts")
	{
		posts.GET("", middleware.OptionalAuthMiddleware(), postController.GetPosts)
		posts.GET("/:id", middleware.OptionalAuthMiddleware(), postController.GetPost)
		posts.GET("/:id/comments", postController.GetComments)

		protected := posts.Group("")
		protected.Use(middleware.AuthMiddleware(rdb))
		{
			protected.POST("", postController.CreatePost)
			protected.PUT("/:id", postController.UpdatePost)
			protected.DELETE("/:id", postController.DeletePost)

			protected.POST("/:id/like", postController.ToggleLike)

			protected.POST("/:id/comments", postController.CreateComment)
			protected.PUT("/:id/comments/:commentId", postController.UpdateComment)
			protected.DELETE("/:id/comments/:commentId", postController.DeleteComment)
		}
	}
}

func SetupRoutes(router *gin.Engine, DB *gorm.DB, rdb *redis.Client, cfg *config.Config) {

	userRepo := repositories.NewUserRepository(DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService, rdb)

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	friendService := &services.FriendService{DB: DB}
	friendController := &controllers.FriendController{Service: friendService}

	uploadService := &services.UploadService{}
	uploadController := &controllers.UploadController{
		Service: uploadService,
	}

	oauthService := services.NewOAuthService(userRepo, rdb, cfg)
	oauthController := controllers.NewOAuthController(oauthService, cfg)

	router.Static("/uploads", "./uploads")
	wsManager := socket.NewWSManager()
	chatHandler := socket.NewChatHandler(wsManager, rdb)

	api := router.Group("/api")
	{
		api.POST("/auth/register", authController.RegisterUser)
		api.POST("/auth/login", authController.LoginUser)
		api.POST("/auth/refresh", authController.RefreshToken)

		api.GET("/auth/oauth/github/login", oauthController.OAuthLogin)
		api.GET("/auth/oauth/github/callback", oauthController.OAuthCallback)

		api.GET("/ws/chat", chatHandler.HandleWS)
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(rdb))
		{
			protected.POST("/auth/logout", authController.LogoutUser)
			protected.GET("users", userController.GetUsers)
			protected.GET("users/:id", userController.GetUser)
			protected.PUT("users/:id", userController.UpdateUser)
			protected.DELETE("users/:id", userController.DeleteUser)

			protected.POST("friends/request/:id", friendController.SendFriendRequest)
			protected.POST("friends/accept/:id", friendController.AcceptFriend)
			protected.POST("friends/follow/:id", friendController.FollowUser)

			protected.POST("upload", uploadController.UploadFile)
		}

		create_post_routes(api, DB, rdb)
	}
}
