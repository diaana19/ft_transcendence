package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// --- Auth Service: LogoutAuthUserService ---

func TestAuthService_LogoutAuthUserService(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "logoutuser", Email: "logout@test.com"}
	db.Create(&user)

	// Generate token
	token, err := utils.GenerateJWT(userID, "logoutuser")
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}

	// Logout
	err = service.LogoutAuthUserService(token, 3600, sharedRDB)
	if err != nil {
		t.Fatalf("LogoutAuthUserService: %v", err)
	}
}

// --- Friend Service: AcceptRequest ---

func TestFriendService_AcceptRequest(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "acceptsender", Email: "as@test.com"}
	receiver := models.User{ID: receiverID, Username: "acceptreceiver", Email: "ar@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	// Send request
	err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	// Accept request
	err = service.AcceptRequest(receiverID, senderID)
	if err != nil {
		t.Fatalf("AcceptRequest: %v", err)
	}

	// Verify friends
	areFriends, _ := service.AreFriends(senderID, receiverID)
	if !areFriends {
		t.Fatal("expected friends after accept")
	}
}

// --- Friend Service: AcceptRequest without pending ---

func TestFriendService_AcceptRequestWithoutPending(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "nopending1", Email: "np1@test.com"}
	user2 := models.User{ID: user2ID, Username: "nopending2", Email: "np2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Accept without pending
	err := service.AcceptRequest(user1ID, user2ID)
	if err == nil {
		t.Fatal("expected error for accept without pending")
	}
}

// --- Post Service: CreateComment with file ---

func TestPostService_CreateCommentWithFile(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "cmtfileuser", Email: "cmtfile@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	// Create file
	fileID := utils.NewID()
	file := models.File{ID: fileID, OwnerID: userID, Path: "/uploads/test.png", Filename: "test.png", MimeType: "image/png", Size: 100, Visibility: "public"}
	db.Create(&file)

	// Create comment with file
	comment, err := service.CreateComment("comment with file", userID, postID, &fileID)
	if err != nil {
		t.Fatalf("CreateComment with file: %v", err)
	}
	if comment.FileID == nil {
		t.Fatal("expected file ID in comment")
	}
}

// --- Post Service: CreateComment on non-existent post ---

func TestPostService_CreateCommentNonExistentPost(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("comment", utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Utils: GenerateJWT ---

func TestUtils_GenerateJWT(t *testing.T) {
	token, err := utils.GenerateJWT("user1", "testuser")
	if err != nil {
		t.Fatalf("GenerateJWT: %v", err)
	}
	if token == "" {
		t.Fatal("expected non-empty token")
	}
}

// --- Utils: ValidateJWT ---

func TestUtils_ValidateJWT(t *testing.T) {
	token, _ := utils.GenerateJWT("user1", "testuser")

	claims, err := utils.ValidateJWT(token)
	if err != nil {
		t.Fatalf("ValidateJWT: %v", err)
	}
	if claims.Subject != "user1" {
		t.Fatalf("expected subject user1, got %s", claims.Subject)
	}
}

// --- Utils: RefreshToken ---

func TestUtils_RefreshToken(t *testing.T) {
	token, _ := utils.GenerateJWT("user1", "testuser")

	refreshed, err := utils.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	if refreshed == "" {
		t.Fatal("expected non-empty refreshed token")
	}
}

// --- Utils: RefreshToken not expired ---

func TestUtils_RefreshTokenNotExpired(t *testing.T) {
	token, _ := utils.GenerateJWT("user1", "testuser")

	// Token should not be refreshed if not near expiration
	refreshed, err := utils.RefreshToken(token)
	if err != nil {
		t.Fatalf("RefreshToken: %v", err)
	}
	// Should return same token if not near expiration
	if refreshed != token {
		t.Fatal("expected same token when not near expiration")
	}
}
