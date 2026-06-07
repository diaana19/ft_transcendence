package test

import (
	"net/http"
	"testing"

	"ft_transcendence/backend/internal/config"
	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/routes"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/socket"
	"ft_transcendence/backend/internal/utils"
)

// ---------------- Repository DB-error branches ----------------

func TestRepositories_DBErrorBranches(t *testing.T) {
	SetupTestEnv()
	bad := brokenDB(t)

	userRepo := repositories.NewUserRepository(bad)
	if _, err := userRepo.Update("u1", models.UpdateUserInput{DisplayName: "x"}); err == nil {
		t.Fatal("userRepo.Update: expected error")
	}
	if err := userRepo.Delete("u1"); err == nil {
		t.Fatal("userRepo.Delete: expected error")
	}

	fileRepo := repositories.NewFileRepository(bad)
	if err := fileRepo.DeleteByOwner("f1", "u1"); err == nil {
		t.Fatal("fileRepo.DeleteByOwner: expected error")
	}

	msgRepo := repositories.NewMessageRepository(bad)
	if _, err := msgRepo.ListConversation("u1", "u2", "", 10); err == nil {
		t.Fatal("ListConversation no-cursor: expected error")
	}
	if _, err := msgRepo.ListConversation("u1", "u2", "cursor", 10); err == nil {
		t.Fatal("ListConversation with-cursor: expected error")
	}
}

// ---------------- Service DB-error branches ----------------

func TestServices_DBErrorBranches(t *testing.T) {
	SetupTestEnv()
	bad := brokenDB(t)

	userSvc := services.NewUserService(repositories.NewUserRepository(bad))
	if err := userSvc.VerifyPassword("u1", "pw"); err == nil {
		t.Fatal("VerifyPassword: expected error")
	}

	twoFASvc := services.NewTwoFAService(repositories.NewUserRepository(bad))
	if _, err := twoFASvc.GenerateSecret("u1"); err == nil {
		t.Fatal("GenerateSecret: expected error storing secret")
	}

	// Friend service: the result.Error (non-RowsAffected) branches.
	fsvc := &services.FriendService{DB: bad}
	if err := fsvc.Unfollow("a", "b"); err == nil {
		t.Fatal("Unfollow on broken db: expected error")
	}
	if err := fsvc.RemoveFriend("a", "b"); err == nil {
		t.Fatal("RemoveFriend on broken db: expected error")
	}
	if err := fsvc.RejectRequest("a", "b"); err == nil {
		t.Fatal("RejectRequest on broken db: expected error")
	}
}

// ---------------- Auth controller: logout service error ----------------

func TestAuthController_LogoutServiceError(t *testing.T) {
	_, db := SetupTestEnv()
	cfg, err := config.Load()
	if err != nil {
		t.Fatalf("config.Load: %v", err)
	}
	// Real DB so JWT validation works, but a closed Redis so blacklisting fails.
	ctrl := routes.Wire(db, brokenRDB(t), cfg)

	token, err := utils.GenerateJWT(utils.NewID(), "logouterr")
	if err != nil {
		t.Fatalf("generate jwt: %v", err)
	}

	c, w := testCtx(http.MethodPost, "/", "")
	c.Request.Header.Set("Authorization", "Bearer "+token)
	ctrl.Auth.LogoutUser(c)
	if w.Code != http.StatusInternalServerError {
		t.Fatalf("LogoutUser service error: expected 500, got %d", w.Code)
	}
}

// ---------------- Chat handler: msgRepo failures ----------------

func TestChatHandler_MsgRepoErrors(t *testing.T) {
	_, db := SetupTestEnv()

	sender := utils.NewID()
	recipient := utils.NewID()
	db.Create(&models.User{ID: sender, Username: "msgerrsender", Email: "mes@test.com"})
	db.Create(&models.User{ID: recipient, Username: "msgerrrecipient", Email: "mer@test.com"})

	manager := socket.NewWSManager()
	notif := services.NewNotificationService(
		repositories.NewNotificationRepositories(db),
		repositories.NewNotificationPubSub(sharedRDB),
	)
	brokenMsg := repositories.NewMessageRepository(brokenDB(t))
	h := socket.NewChatHandler(
		manager, sharedRDB, notif, brokenMsg,
		repositories.NewUserRepository(db), repositories.NewFileRepository(db), "",
	)

	client := &socket.Client{ID: sender, Username: "msgerrsender", Send: make(chan []byte, 256)}

	// handleOpen -> ListConversation fails -> history load error branch.
	h.HandleMessage(client, []byte(`{"action":"open","peer_id":"`+recipient+`"}`))
	// handleDM -> recipient lookup ok, msgRepo.Create fails -> save error branch.
	h.HandleMessage(client, []byte(`{"action":"message","content":"hi","recipient_id":"`+recipient+`"}`))

	for {
		select {
		case <-client.Send:
		default:
			return
		}
	}
}
