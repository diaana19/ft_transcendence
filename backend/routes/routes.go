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
	notifRepo := repositories.NewNotificationRepositories(DB)
	notifPubSub := repositories.NewNotiticationPubSub(rdb)
	notifService := services.NewNotificationService(notifRepo, notifPubSub)

	postRepo := repositories.NewPostRepository(DB)
	postService := services.NewPostService(postRepo)
	postController := controllers.NewPostController(postService, notifService)

	posts := api.Group("/posts")
	{
		posts.GET("", middleware.OptionalAuthMiddleware(), postController.GetPosts)
		posts.GET("/:id", middleware.OptionalAuthMiddleware(), postController.GetPost)
		posts.GET("/user/:userId", middleware.OptionalAuthMiddleware(), postController.GetPostsByUser)
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
	notifRepo := repositories.NewNotificationRepositories(DB)
	notifPubSub := repositories.NewNotiticationPubSub(rdb)
	notifService := services.NewNotificationService(notifRepo, notifPubSub)
	notifController := controllers.NewNotificationController(notifService)

	msgRepo := repositories.NewMessageRepository(DB)
	msgController := controllers.NewMsgController(msgRepo)

	userRepo := repositories.NewUserRepository(DB)
	authService := services.NewAuthService(userRepo)
	authController := controllers.NewAuthController(authService, rdb)

	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	friendService := &services.FriendService{DB: DB}
	friendController := &controllers.FriendController{Service: friendService, NotificationService: notifService}

	uploadService := &services.UploadService{}
	uploadController := &controllers.UploadController{Service: uploadService}

	wsManager := socket.NewWSManager()
	chatHandler := socket.NewChatHandler(wsManager, rdb, notifService, msgRepo)
	oauthService := services.NewOAuthService(userRepo, rdb, cfg)
	oauthController := controllers.NewOAuthController(oauthService, cfg)

	chatService := services.NewChatService(msgRepo, userRepo)
	chatController := controllers.NewChatController(chatService)

	router.Static("/uploads", "./uploads")

	api := router.Group("/api")
	{
		api.POST("/auth/register", authController.RegisterUser)
		api.POST("/auth/login", middleware.RateLimitMiddleware(rdb), authController.LoginUser)
		api.POST("/auth/refresh", middleware.RateLimitMiddleware(rdb), authController.RefreshToken)
		api.GET("/ws/chat", middleware.WSAuthMiddleware(), chatHandler.HandleWS)

		api.GET("/auth/oauth/github/login", oauthController.OAuthLogin)
		api.GET("/auth/oauth/github/callback", oauthController.OAuthCallback)

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

			protected.POST("notification", notifController.GetUnread)
			protected.PUT("notification/read", notifController.MarkAllRead)

			// Room message history — served by the WebSocket primary path
			protected.GET("rooms/:roomId/messages", msgController.GetRoomMsg)
			protected.GET("messages/:messageId/replies", msgController.GetReplies)

			// Poll/DM fallback — used when WebSocket is unavailable
			protected.POST("chat/messages", chatController.SendMessage)
			protected.GET("chat/messages", chatController.ListConversation)
			protected.GET("chat/poll", chatController.Poll)

			protected.POST("upload", uploadController.UploadFile)
		}

		create_post_routes(api, DB, rdb)
	}
}
