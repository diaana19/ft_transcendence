package routes

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/controllers"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/socket"
)

// Controllers aggregates every HTTP handler the API exposes plus the
// WebSocket chat handler. Built once at startup by Wire and passed to
// SetupRoutes for registration.
type Controllers struct {
	Auth         *controllers.AuthController
	TwoFA        *controllers.TwoFAController
	User         *controllers.UserController
	Friend       *controllers.FriendController
	Post         *controllers.PostController
	Notification *controllers.NotificationController
	Msg          *controllers.MsgController
	Chat         *controllers.ChatController
	Upload       *controllers.UploadController
	GDPR         *controllers.GDPRController
	OAuth        *controllers.OAuthController
	Search       *controllers.SearchController
	Gamification *controllers.GamificationController
	ChatWS       *socket.ChatHandler
}

// Wire builds the full dependency graph (repositories -> services -> controllers)
// from the shared infrastructure clients. Construction order respects the
// dependency direction: every value used by a later constructor is declared
// before it.
func Wire(pdb *gorm.DB, rdb *redis.Client, cfg *config.Config) *Controllers {
	userRepo := repositories.NewUserRepository(pdb)
	msgRepo := repositories.NewMessageRepository(pdb)
	postRepo := repositories.NewPostRepository(pdb)
	fileRepo := repositories.NewFileRepository(pdb)
	notifRepo := repositories.NewNotificationRepositories(pdb)
	notifPubSub := repositories.NewNotificationPubSub(rdb)

	authService := services.NewAuthService(userRepo)
	twoFAService := services.NewTwoFAService(userRepo)
	userService := services.NewUserService(userRepo)
	friendService := &services.FriendService{DB: pdb}
	postService := services.NewPostService(postRepo)
	notifService := services.NewNotificationService(notifRepo, notifPubSub)
	uploadService := services.NewUploadService(fileRepo)
	gdprService := services.NewGDPRService(pdb)
	chatService := services.NewChatService(msgRepo, userRepo)
	oauthService := services.NewOAuthService(userRepo, rdb, cfg)
	searchService := services.NewSearchService(userRepo, msgRepo, postRepo)
	gamificationService := services.NewGamificationService(pdb, friendService)

	wsManager := socket.NewWSManager()
	chatWS := socket.NewChatHandler(wsManager, rdb, notifService, msgRepo, fileRepo, cfg.FrontendURL)

	return &Controllers{
		Auth:         controllers.NewAuthController(authService, twoFAService, rdb),
		TwoFA:        controllers.NewTwoFAController(twoFAService),
		User:         controllers.NewUserController(userService, friendService),
		Friend:       &controllers.FriendController{Service: friendService, NotificationService: notifService},
		Post:         controllers.NewPostController(postService, notifService, uploadService),
		Notification: controllers.NewNotificationController(notifService),
		Msg:          controllers.NewMsgController(msgRepo),
		Chat:         controllers.NewChatController(chatService),
		Upload:       &controllers.UploadController{Service: uploadService, FriendService: friendService},
		GDPR:         controllers.NewGDPRController(gdprService),
		OAuth:        controllers.NewOAuthController(oauthService, cfg),
		Search:       controllers.NewSearchController(searchService),
		Gamification: controllers.NewGamificationController(gamificationService),
		ChatWS:       chatWS,
	}
}
