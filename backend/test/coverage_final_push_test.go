package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

// --- Post Service: CreateComment ---

func TestPostService_CreateCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "createcommentfinal", Email: "createcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	// Create comment
	comment, err := service.CreateComment("new comment", userID, postID, nil)
	if err != nil {
		t.Fatalf("CreateComment: %v", err)
	}
	if comment.Content != "new comment" {
		t.Fatalf("expected 'new comment', got %q", comment.Content)
	}
}

// --- Post Service: CreateComment empty ---

func TestPostService_CreateCommentEmptyFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("", utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for empty comment")
	}
}

// --- Post Service: CreateComment too long ---

func TestPostService_CreateCommentTooLongFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	longContent := ""
	for i := 0; i < 300; i++ {
		longContent += "a"
	}

	_, err := service.CreateComment(longContent, utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for comment over 280 chars")
	}
}

// --- Post Service: CreateComment not found ---

func TestPostService_CreateCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("comment", utils.NewID(), utils.NewID(), nil)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Post Service: UpdateComment ---

func TestPostService_UpdateCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updcommentfinal", Email: "updcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("original", userID, postID, nil)

	// Update comment
	updated, err := service.UpdateComment(comment.ID, models.UpdateCommentInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

// --- Post Service: UpdateComment not found ---

func TestPostService_UpdateCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.UpdateComment(utils.NewID(), models.UpdateCommentInput{Content: "updated"}, utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

// --- Post Service: UpdateComment forbidden ---

func TestPostService_UpdateCommentForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "updcforbiddenfinal", Email: "updcforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherupdcfinal", Email: "otherupdcfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	comment, _ := service.CreateComment("comment", userID, postID, nil)

	// Try to update other user's comment
	_, err := service.UpdateComment(comment.ID, models.UpdateCommentInput{Content: "hacked"}, otherUserID)
	if err == nil {
		t.Fatal("expected error for updating other user's comment")
	}
}

// --- Post Service: DeleteComment ---

func TestPostService_DeleteCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delcommentfinal", Email: "delcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("to delete", userID, postID, nil)

	// Delete comment
	err := service.DeleteComment(comment.ID, userID)
	if err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
}

// --- Post Service: DeleteComment not found ---

func TestPostService_DeleteCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeleteComment(utils.NewID(), utils.NewID())
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

// --- Post Service: DeleteComment forbidden ---

func TestPostService_DeleteCommentForbiddenFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	otherUserID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "delcforbiddenfinal", Email: "delcforbiddenfinal@test.com"}
	otherUser := models.User{ID: otherUserID, Username: "otherdelcfinal", Email: "otherdelcfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&otherUser)
	repo.Create(&post)

	comment, _ := service.CreateComment("comment", userID, postID, nil)

	// Try to delete other user's comment
	err := service.DeleteComment(comment.ID, otherUserID)
	if err == nil {
		t.Fatal("expected error for deleting other user's comment")
	}
}

// --- Post Service: ReactToPost ---

func TestPostService_ReactToPostFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "reactpostfinal", Email: "reactpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	// Like
	value, _, err := service.ReactToPost(userID, postID, 1)
	if err != nil {
		t.Fatalf("ReactToPost like: %v", err)
	}
	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}

	// Unlike (toggle)
	value, _, err = service.ReactToPost(userID, postID, 1)
	if err != nil {
		t.Fatalf("ReactToPost unlike: %v", err)
	}
	if value != 0 {
		t.Fatalf("expected 0, got %d", value)
	}
}

// --- Post Service: ReactToPost not found ---

func TestPostService_ReactToPostNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToPost(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent post")
	}
}

// --- Post Service: ReactToComment ---

func TestPostService_ReactToCommentFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "reactcommentfinal", Email: "reactcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	repo.Create(&post)
	repo.CreateComment(&comment)

	// Like comment
	value, _, err := service.ReactToComment(userID, commentID, 1)
	if err != nil {
		t.Fatalf("ReactToComment: %v", err)
	}
	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}
}

// --- Post Service: ReactToComment not found ---

func TestPostService_ReactToCommentNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, _, err := service.ReactToComment(utils.NewID(), utils.NewID(), 1)
	if err == nil {
		t.Fatal("expected error for non-existent comment")
	}
}

// --- Post Service: GetPostReaction ---

func TestPostService_GetPostReactionFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "getreactpostfinal", Email: "getreactpostfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	// Set reaction
	repo.SetPostReaction(userID, postID, 1)

	// Get reaction
	val, err := service.GetPostReaction(userID, postID)
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

// --- Post Service: GetCommentReaction ---

func TestPostService_GetCommentReactionFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	// Create test user, post, and comment
	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "getreactcommentfinal", Email: "getreactcommentfinal@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	repo.Create(&post)
	repo.CreateComment(&comment)

	// Set reaction
	repo.SetReplyReaction(userID, commentID, 1)

	// Get reaction
	val, err := service.GetCommentReaction(userID, commentID)
	if err != nil {
		t.Fatalf("GetCommentReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

// --- Auth Service: LoginAuthUserService ---

func TestAuthService_LoginAuthUserServiceFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create test user with password
	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginfinal", Email: "loginfinal@test.com", Password: &hash}
	db.Create(&user)

	// Login with email
	_, err := service.LoginAuthUserService("loginfinal@test.com", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with email: %v", err)
	}

	// Login with username
	_, err = service.LoginAuthUserService("loginfinal", password)
	if err != nil {
		t.Fatalf("LoginAuthUserService with username: %v", err)
	}
}

// --- Auth Service: LoginAuthUserService wrong password ---

func TestAuthService_LoginAuthUserServiceWrongPasswordFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create test user with password
	userID := utils.NewID()
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := models.User{ID: userID, Username: "loginwpfinal", Email: "loginwpfinal@test.com", Password: &hash}
	db.Create(&user)

	// Login with wrong password
	_, err := service.LoginAuthUserService("loginwpfinal@test.com", "wrongpassword")
	if err == nil {
		t.Fatal("expected error for wrong password")
	}
}

// --- Auth Service: LoginAuthUserService not found ---

func TestAuthService_LoginAuthUserServiceNotFoundFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Login with non-existent user
	_, err := service.LoginAuthUserService("nonexistent@test.com", "password")
	if err == nil {
		t.Fatal("expected error for non-existent user")
	}
}

// --- Auth Service: CreateAuthUserService ---

func TestAuthService_CreateAuthUserServiceFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create new user
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user := &models.User{
		Username: "createfinal",
		Email:    "createfinal@test.com",
		Password: &hash,
	}
	created, err := service.CreateAuthUserService(user)
	if err != nil {
		t.Fatalf("CreateAuthUserService: %v", err)
	}
	if created.Username != "createfinal" {
		t.Fatalf("expected username 'createfinal', got %q", created.Username)
	}
}

// --- Auth Service: CreateAuthUserService duplicate ---

func TestAuthService_CreateAuthUserServiceDuplicateFinal(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewAuthService(repo)

	// Create first user
	password := "StrongPass123!"
	hash, _ := utils.HashString(password)
	user1 := &models.User{
		Username: "dupfinal",
		Email:    "dupfinal@test.com",
		Password: &hash,
	}
	_, err := service.CreateAuthUserService(user1)
	if err != nil {
		t.Fatalf("CreateAuthUserService first: %v", err)
	}

	// Create duplicate username
	user2 := &models.User{
		Username: "dupfinal",
		Email:    "dupfinal2@test.com",
		Password: &hash,
	}
	_, err = service.CreateAuthUserService(user2)
	if err == nil {
		t.Fatal("expected error for duplicate username")
	}
}

// --- Friend Service: AcceptRequest ---

func TestFriendService_AcceptRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "acceptsenderfinal", Email: "asfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "acceptreceiverfinal", Email: "arfinal@test.com"}
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

func TestFriendService_AcceptRequestWithoutPendingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "nopendingfinal1", Email: "npfinal1@test.com"}
	user2 := models.User{ID: user2ID, Username: "nopendingfinal2", Email: "npfinal2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Accept without pending
	err := service.AcceptRequest(user1ID, user2ID)
	if err == nil {
		t.Fatal("expected error for accept without pending")
	}
}

// --- Friend Service: RejectRequest ---

func TestFriendService_RejectRequestFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	senderID := utils.NewID()
	receiverID := utils.NewID()
	sender := models.User{ID: senderID, Username: "rejectsenderfinal", Email: "rsfinal@test.com"}
	receiver := models.User{ID: receiverID, Username: "rejectreceiverfinal", Email: "rrfinal@test.com"}
	db.Create(&sender)
	db.Create(&receiver)

	// Send request
	err := service.SendRequest(senderID, receiverID)
	if err != nil {
		t.Fatalf("SendRequest: %v", err)
	}

	// Reject request
	err = service.RejectRequest(receiverID, senderID)
	if err != nil {
		t.Fatalf("RejectRequest: %v", err)
	}

	// Verify not friends
	areFriends, _ := service.AreFriends(senderID, receiverID)
	if areFriends {
		t.Fatal("expected not friends after reject")
	}
}

// --- Friend Service: RejectRequest without pending ---

func TestFriendService_RejectRequestWithoutPendingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "norejectfinal1", Email: "nr1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "norejectfinal2", Email: "nr2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Reject without pending
	err := service.RejectRequest(user1ID, user2ID)
	if err == nil {
		t.Fatal("expected error for reject without pending")
	}
}

// --- Friend Service: GetFollowers ---

func TestFriendService_GetFollowersFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "follower1final", Email: "f1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "follower2final", Email: "f2final@test.com"}
	target := models.User{ID: targetID, Username: "targetfinal", Email: "targetfinal@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&target)

	// Create follows
	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	// Get followers
	followers, err := service.GetFollowers(targetID)
	if err != nil {
		t.Fatalf("GetFollowers: %v", err)
	}
	if len(followers) != 2 {
		t.Fatalf("expected 2 followers, got %d", len(followers))
	}
}

// --- Friend Service: GetFollowing ---

func TestFriendService_GetFollowingFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	userID := utils.NewID()
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user := models.User{ID: userID, Username: "followinguserfinal", Email: "fufinal@test.com"}
	user1 := models.User{ID: user1ID, Username: "following1final", Email: "f1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "following2final", Email: "f2final@test.com"}
	db.Create(&user)
	db.Create(&user1)
	db.Create(&user2)

	// Create follows
	service.Follow(userID, user1ID)
	service.Follow(userID, user2ID)

	// Get following
	following, err := service.GetFollowing(userID)
	if err != nil {
		t.Fatalf("GetFollowing: %v", err)
	}
	if len(following) != 2 {
		t.Fatalf("expected 2 following, got %d", len(following))
	}
}

// --- Friend Service: GetFriends ---

func TestFriendService_GetFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "friend1final", Email: "fr1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "friend2final", Email: "fr2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Make friends
	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	// Get friends
	friends, err := service.GetFriends(user1ID)
	if err != nil {
		t.Fatalf("GetFriends: %v", err)
	}
	if len(friends) != 1 {
		t.Fatalf("expected 1 friend, got %d", len(friends))
	}
}

// --- Friend Service: AreFriends ---

func TestFriendService_AreFriendsFinal(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "arefriends1final", Email: "af1final@test.com"}
	user2 := models.User{ID: user2ID, Username: "arefriends2final", Email: "af2final@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Not friends initially
	areFriends, err := service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends: %v", err)
	}
	if areFriends {
		t.Fatal("expected not friends initially")
	}

	// Make friends
	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	// Now friends
	areFriends, err = service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends after accept: %v", err)
	}
	if !areFriends {
		t.Fatal("expected friends after accept")
	}
}
