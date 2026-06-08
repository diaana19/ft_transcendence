package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"

	"ft_transcendence/backend/internal/controllers"
	"ft_transcendence/backend/internal/middleware"
	"ft_transcendence/backend/internal/socket"
)

func SetupRoutes(router *gin.Engine, c *Controllers, rdb *redis.Client) {
	api := router.Group("/api")

	registerPublicAuthRoutes(api, c.Auth, rdb)
	registerOAuthRoutes(api, c.OAuth)
	registerWebSocketRoutes(api, c.ChatWS)
	registerPostRoutes(api, rdb, c.Post)

	api.GET("/trends", c.Post.GetTrends)

	api.GET("/users/:id/online", onlineStatusHandler(c.WSManager))

	protected := api.Group("", middleware.AuthMiddleware(rdb))
	registerProtectedAuthRoutes(protected, c.Auth)
	registerUserRoutes(protected, c.User)
	registerFriendRoutes(protected, c.Friend)
	registerNotificationRoutes(protected, c.Notification)
	registerUploadRoutes(api, rdb, c.Upload)
	registerGDPRRoutes(protected, c.GDPR)
	registerTwoFARoutes(protected, c.TwoFA)
	registerGamificationRoutes(protected, c.Gamification)
}

// onlineStatusHandler godoc
// @Summary   Check whether a user is currently online (has an active WebSocket connection)
// @Tags      users
// @Produce   json
// @Param     id path string true "user id"
// @Success   200 {object} map[string]bool
// @Router    /users/{id}/online [get]
func onlineStatusHandler(wsManager *socket.WSManager) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"online": wsManager.IsOnline(ctx.Param("id"))})
	}
}

func registerPublicAuthRoutes(api *gin.RouterGroup, c *controllers.AuthController, rdb *redis.Client) {
	auth := api.Group("/auth", middleware.RateLimitMiddleware(rdb))
	auth.POST("/register", c.RegisterUser)
	auth.POST("/login", c.LoginUser)
	auth.POST("/refresh", c.RefreshToken)
	auth.POST("/2fa/verify", c.Verify2FA)
	auth.POST("/forgot-password", c.ForgotPassword)
}

func registerOAuthRoutes(api *gin.RouterGroup, c *controllers.OAuthController) {
	oauth := api.Group("/auth/oauth/github")
	oauth.GET("/login", c.OAuthLogin)
	oauth.GET("/callback", c.OAuthCallback)
}

func registerWebSocketRoutes(api *gin.RouterGroup, c *socket.ChatHandler) {
	api.GET("/ws/chat", middleware.WSAuthMiddleware(), c.HandleWS)
}

func registerProtectedAuthRoutes(protected *gin.RouterGroup, c *controllers.AuthController) {
	protected.POST("/auth/logout", c.LogoutUser)
	protected.GET("/auth/me", c.Me)
}

func registerUserRoutes(protected *gin.RouterGroup, c *controllers.UserController) {
	protected.GET("/users", c.GetUsers)
	protected.GET("/users/:id", c.GetUser)
	protected.PUT("/users/:id", c.UpdateUser)
	protected.DELETE("/users/:id", c.DeleteUser)
}

func registerFriendRoutes(protected *gin.RouterGroup, c *controllers.FriendController) {
	protected.POST("/friends/request/:id", c.SendFriendRequest)
	protected.POST("/friends/accept/:id", c.AcceptFriend)
	protected.POST("/friends/reject/:id", c.RejectFriendRequest)
	protected.DELETE("/friends/:id", c.RemoveFriend)
	protected.POST("/friends/follow/:id", c.FollowUser)
	protected.DELETE("/friends/follow/:id", c.UnfollowUser)
	protected.GET("/users/:id/followers", c.GetFollowers)
	protected.GET("/users/:id/following", c.GetFollowing)
	protected.GET("/users/:id/friends", c.GetFriends)
}

func registerNotificationRoutes(protected *gin.RouterGroup, c *controllers.NotificationController) {
	protected.GET("/notification", c.GetUnread)
	protected.PATCH("/notification/read", c.MarkAllRead)
	protected.PATCH("/notification/:id/read", c.MarkRead)
}

func registerUploadRoutes(api *gin.RouterGroup, rdb *redis.Client, c *controllers.UploadController) {
	api.GET("/files/:id", middleware.OptionalAuthMiddleware(rdb), c.ServeFile)

	protected := api.Group("", middleware.AuthMiddleware(rdb))
	protected.POST("/upload", c.UploadFile)
}

func registerGDPRRoutes(protected *gin.RouterGroup, c *controllers.GDPRController) {
	protected.GET("/gdpr/export", c.ExportUserData)
	protected.DELETE("/gdpr/delete", c.DeleteUserData)
}

func registerTwoFARoutes(protected *gin.RouterGroup, c *controllers.TwoFAController) {
	protected.POST("/2fa/setup", c.Setup)
	protected.POST("/2fa/enable", c.Enable)
	protected.POST("/2fa/disable", c.Disable)
}

func registerGamificationRoutes(protected *gin.RouterGroup, c *controllers.GamificationController) {
	protected.GET("/users/:id/gamification", c.Get)
	protected.GET("/leaderboard", c.GetLeaderboard)
}

func registerPostRoutes(api *gin.RouterGroup, rdb *redis.Client, c *controllers.PostController) {
	posts := api.Group("/posts")
	posts.GET("", middleware.OptionalAuthMiddleware(rdb), c.GetPosts)
	posts.GET("/user/:userId", middleware.OptionalAuthMiddleware(rdb), c.GetPostsByUser)
	posts.GET("/:id", middleware.OptionalAuthMiddleware(rdb), c.GetPost)
	posts.GET("/:id/comments", middleware.OptionalAuthMiddleware(rdb), c.GetComments)

	protected := posts.Group("", middleware.AuthMiddleware(rdb))
	protected.POST("", c.CreatePost)
	protected.PUT("/:id", c.UpdatePost)
	protected.DELETE("/:id", c.DeletePost)
	protected.POST("/:id/react", c.React)
	protected.POST("/:id/comments", c.CreateComment)
	protected.PUT("/:id/comments/:commentId", c.UpdateComment)
	protected.DELETE("/:id/comments/:commentId", c.DeleteComment)
	protected.POST("/:id/comments/:commentId/react", c.ReactComment)
}
