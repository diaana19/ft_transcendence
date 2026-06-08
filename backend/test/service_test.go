package test

import (
	"testing"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/services"
	"ft_transcendence/backend/internal/utils"
)

func TestPostService_GetPostsByTag(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "svctaguser", Email: "svctag@test.com"}
	db.Create(&user)

	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post with #go", Tags: []string{"#go"}}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post with #rust", Tags: []string{"#rust"}}
	repo.Create(&post1)
	repo.Create(&post2)

	posts, total, err := service.GetPostsByTag("#go", 10, 0)
	if err != nil {
		t.Fatalf("GetPostsByTag: %v", err)
	}
	if len(posts) != 1 {
		t.Fatalf("expected 1 post, got %d", len(posts))
	}
	if total != 1 {
		t.Fatalf("expected total 1, got %d", total)
	}
}

func TestPostService_CreateComment(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "svccommentuser", Email: "svccomment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, err := service.CreateComment("nice post", userID, postID, nil)
	if err != nil {
		t.Fatalf("CreateComment: %v", err)
	}
	if comment.Content != "nice post" {
		t.Fatalf("expected 'nice post', got %q", comment.Content)
	}
}

func TestPostService_CreateCommentEmptyContent(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	_, err := service.CreateComment("", "user", "post", nil)
	if err == nil {
		t.Fatal("expected error for empty content")
	}
}

func TestPostService_CreateCommentOverLimit(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	longContent := ""
	for i := 0; i < 300; i++ {
		longContent += "a"
	}

	_, err := service.CreateComment(longContent, "user", "post", nil)
	if err == nil {
		t.Fatal("expected error for content over 280 chars")
	}
}

func TestPostService_GetComments(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "svcgetcommentuser", Email: "svcgetcomment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	service.CreateComment("comment1", userID, postID, nil)
	service.CreateComment("comment2", userID, postID, nil)

	comments, err := service.GetComments(postID)
	if err != nil {
		t.Fatalf("GetComments: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestPostService_UpdateComment(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "svcupdcommentuser", Email: "svcupdcomment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("original", userID, postID, nil)

	updated, err := service.UpdateComment(comment.ID, models.UpdateCommentInput{Content: "updated"}, userID)
	if err != nil {
		t.Fatalf("UpdateComment: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

func TestPostService_DeleteComment(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "svcdelcommentuser", Email: "svcdelcomment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	repo.Create(&post)

	comment, _ := service.CreateComment("to delete", userID, postID, nil)

	err := service.DeleteComment(comment.ID, userID)
	if err != nil {
		t.Fatalf("DeleteComment: %v", err)
	}
}

func TestPostService_DeleteCommentNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	err := service.DeleteComment(utils.NewID(), "user")
	if err == nil {
		t.Fatal("expected error for nonexistent comment")
	}
}

func TestPostService_ReactToPost(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "svcreactuser", Email: "svcreact@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	repo.Create(&post)

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

func TestPostService_ReactToComment(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)
	service := services.NewPostService(repo)

	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "svcreactcommentuser", Email: "svcreactcomment@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: utils.NewID(), Content: "comment"}
	db.Create(&user)
	repo.Create(&post)
	repo.CreateComment(&comment)

	value, _, err := service.ReactToComment(userID, commentID, 1)
	if err != nil {
		t.Fatalf("ReactToComment: %v", err)
	}
	if value != 1 {
		t.Fatalf("expected 1, got %d", value)
	}
}

func TestUserService_GetUserByUsername(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	userID := utils.NewID()
	user := models.User{ID: userID, Username: "svcusernameuser", Email: "svcusername@test.com"}
	db.Create(&user)

	got, err := service.GetUserByUsername("svcusernameuser")
	if err != nil {
		t.Fatalf("GetUserByUsername: %v", err)
	}
	if got.ID != userID {
		t.Fatalf("expected ID %s, got %s", userID, got.ID)
	}
}

func TestUserService_GetUserByUsernameNotFound(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewUserRepository(db)
	service := services.NewUserService(repo)

	_, err := service.GetUserByUsername("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent username")
	}
}

func TestFriendService_CountFollowers(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	targetID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "countfollower1", Email: "count1@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollower2", Email: "count2@test.com"}
	user3 := models.User{ID: targetID, Username: "counttarget", Email: "counttarget@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&user3)

	service.Follow(user1ID, targetID)
	service.Follow(user2ID, targetID)

	count, err := service.CountFollowers(targetID)
	if err != nil {
		t.Fatalf("CountFollowers: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 followers, got %d", count)
	}
}

func TestFriendService_CountFollowing(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user3ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "countfollowing1", Email: "countf1@test.com"}
	user2 := models.User{ID: user2ID, Username: "countfollowing2", Email: "countf2@test.com"}
	user3 := models.User{ID: user3ID, Username: "countfollowing3", Email: "countf3@test.com"}
	db.Create(&user1)
	db.Create(&user2)
	db.Create(&user3)

	service.Follow(user1ID, user2ID)
	service.Follow(user1ID, user3ID)

	count, err := service.CountFollowing(user1ID)
	if err != nil {
		t.Fatalf("CountFollowing: %v", err)
	}
	if count != 2 {
		t.Fatalf("expected 2 following, got %d", count)
	}
}

func TestFriendService_AreFriends(t *testing.T) {
	_, db := SetupTestEnv()
	service := &services.FriendService{DB: db}

	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "friendcheck1", Email: "fc1@test.com"}
	user2 := models.User{ID: user2ID, Username: "friendcheck2", Email: "fc2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	areFriends, err := service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends: %v", err)
	}
	if areFriends {
		t.Fatal("expected not friends initially")
	}

	service.SendRequest(user1ID, user2ID)
	service.AcceptRequest(user2ID, user1ID)

	areFriends, err = service.AreFriends(user1ID, user2ID)
	if err != nil {
		t.Fatalf("AreFriends after accept: %v", err)
	}
	if !areFriends {
		t.Fatal("expected friends after accept")
	}
}
