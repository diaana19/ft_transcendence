package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestPostService_UpdatePost(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updpostuser", Email: "updpost@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "original"}
	db.Create(&user)
	repo.Create(&post)

	updated, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdatePost: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

func TestPostService_UpdatePostNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.UpdatePost(utils.NewID(), models.UpdatePostInput{Content: "updated"}, utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_UpdatePostForbidden(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "forbiddenuser", Email: "forbidden@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otheruser", Email: "other@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	_, err := service.UpdatePost(postID, models.UpdatePostInput{Content: "hacked"}, otherUserID)
	if err == nil {
		t.Fatal("expected error for updating other user's post")
	}
}

func TestPostService_DeletePost(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delpostuser", Email: "delpost@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "to be deleted"}
	db.Create(&user)
	repo.Create(&post)

	err := service.DeletePost(postID, userID)
	if err != nil {
		t.Fatalf("DeletePost: %v", err)
	}

	_, err = repo.GetByID(postID)
	if err == nil {
		t.Fatal("expected error for deleted post")
	}
}

func TestPostService_DeletePostNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeletePost(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_DeletePostForbidden(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delforbiddenuser", Email: "delforbidden@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherdeluser", Email: "otherdel@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	err := service.DeletePost(postID, otherUserID)
	if err == nil {
		t.Fatal("expected error for deleting other user's post")
	}
}

func TestPostService_ReactToPostErrorPath(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToPost(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

func TestPostService_ReactToCommentErrorPath(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToComment(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

func TestPostService_GetPostReaction(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "getreactuser", Email: "getreact@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	repo.SetPostReaction(userID, postID, 1)

	val, err := service.GetPostReaction(userID, postID)
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestPostService_GetCommentReaction(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "getcommentreactuser", Email: "getcommentreact@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	db.Create(&post)
	repo.CreateComment(&comment)

	repo.SetReplyReaction(userID, commentID, 1)

	val, err := service.GetCommentReaction(userID, commentID)
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestUserService_GetUsers(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "getusers1", Email: "gu1@test.com"}
	user2 := models.User{ID: user2ID, Username: "getusers2", Email: "gu2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	users, err := service.GetUsers()
	if err != nil {
		t.Fatalf("GetUsers: %v", err)
	}
	if len(users) < 2 {
		t.Fatalf("expected at least 2 users, got %d", len(users))
	}
}

func TestAuthService_LoginAuthUserService(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginuser", Email: "login@test.com", Password: &hash}
	db.Create(&user)

	_, err := service.LoginAuthUserService("login@test.com", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with email: %v", err)
	}

	_, err = service.LoginAuthUserService("loginuser", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with username: %v", err)
	}
}

func TestAuthService_LoginAuthUserServiceWrongPassword(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginwpuser", Email: "loginwp@test.com", Password: &hash}
	db.Create(&user)

	_, err := service.LoginAuthUserService("loginwp@test.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

func TestAuthService_LoginAuthUserServiceNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	_, err := service.LoginAuthUserService("nonexistent@test.com", "password")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

func TestAuthService_CreateAuthUserService(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := &models.User{
		Username: "newuser",
		Email:    "new@test.com",
		Password: &hash,
	}
	created, err := service.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("CreateAuthUserService: %v", err)
	}
	if created.Username != "newuser" {
		t.Fatalf("expected username 'newuser', got %q", created.Username)
	}
}

func TestAuthService_CreateAuthUserServiceDuplicate(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user1 := &models.User{
		Username: "dupuser",
		Email:    "dup@test.com",
		Password: &hash,
	}
	_, err := service.CreateAuthUserService(user1)
	if err != nil {
		t.Fatalf("CreateAuthUserService first: %v", err)
	}

	user2 := &models.User{
		Username: "dupuser",
		Email:    "dup2@test.com",
		Password: &hash,
	}
	_, err = service.CreateAuthUserService(user2)
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}
