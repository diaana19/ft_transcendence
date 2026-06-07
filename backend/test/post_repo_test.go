package test

import (
	"testing"
	"time"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

func TestPostRepository_GetByTag(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "taguser", Email: "tag@test.com"}
	db.Create(&user)

	// Create posts with tags
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post with #golang", Tags: []string{"#golang"}}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post with #rust", Tags: []string{"#rust"}}
	post3 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post with both", Tags: []string{"#golang", "#rust"}}
	repo.Create(&post1)
	repo.Create(&post2)
	repo.Create(&post3)

	// Get by tag
	posts, total, err := repo.GetByTag("#golang", 10, 0)
	if err != nil {
		t.Fatalf("GetByTag: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts with #golang, got %d", len(posts))
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
}

func TestPostRepository_GetByTagEmpty(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	posts, total, err := repo.GetByTag("#nonexistent", 10, 0)
	if err != nil {
		t.Fatalf("GetByTag empty: %v", err)
	}
	if len(posts) != 0 {
		t.Fatalf("expected 0 posts, got %d", len(posts))
	}
	if total != 0 {
		t.Fatalf("expected total 0, got %d", total)
	}
}

func TestPostRepository_TopTags(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "trenduser", Email: "trend@test.com"}
	db.Create(&user)

	// Create posts with tags
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1", Tags: []string{"#golang"}}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2", Tags: []string{"#golang"}}
	post3 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post3", Tags: []string{"#rust"}}
	repo.Create(&post1)
	repo.Create(&post2)
	repo.Create(&post3)

	// Get top tags
	since := time.Now().Add(-24 * time.Hour)
	tags, err := repo.TopTags(since, 10)
	if err != nil {
		t.Fatalf("TopTags: %v", err)
	}
	if len(tags) < 2 {
		t.Fatalf("expected at least 2 tags, got %d", len(tags))
	}

	// First tag should be #golang (most used)
	if tags[0].Tag != "#golang" {
		t.Fatalf("expected #golang as top tag, got %s", tags[0].Tag)
	}
}

func TestPostRepository_GetAll(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "alluser", Email: "all@test.com"}
	db.Create(&user)

	// Create posts
	post1 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: userID, Content: "post2"}
	repo.Create(&post1)
	repo.Create(&post2)

	// Get all
	posts, total, err := repo.GetAll(10, 0)
	if err != nil {
		t.Fatalf("GetAll: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d", len(posts))
	}
	if total != 2 {
		t.Fatalf("expected total 2, got %d", total)
	}
}

func TestPostRepository_GetByAuthorID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test users
	user1ID := utils.NewID()
	user2ID := utils.NewID()
	user1 := models.User{ID: user1ID, Username: "author1", Email: "author1@test.com"}
	user2 := models.User{ID: user2ID, Username: "author2", Email: "author2@test.com"}
	db.Create(&user1)
	db.Create(&user2)

	// Create posts
	post1 := models.Post{ID: utils.NewID(), AuthorID: user1ID, Content: "post by author1"}
	post2 := models.Post{ID: utils.NewID(), AuthorID: user1ID, Content: "another by author1"}
	post3 := models.Post{ID: utils.NewID(), AuthorID: user2ID, Content: "post by author2"}
	repo.Create(&post1)
	repo.Create(&post2)
	repo.Create(&post3)

	// Get by author
	posts, err := repo.GetByAuthorID(user1ID)
	if err != nil {
		t.Fatalf("GetByAuthorID: %v", err)
	}
	if len(posts) != 2 {
		t.Fatalf("expected 2 posts by author1, got %d", len(posts))
	}
}

func TestPostRepository_Update(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "updateuser", Email: "update@test.com"}
	db.Create(&user)

	// Create post
	postID := utils.NewID()
	post := models.Post{ID: postID, AuthorID: userID, Content: "original"}
	repo.Create(&post)

	// Update
	updated, err := repo.Update(postID, models.UpdatePostInput{Content: "updated"})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Content != "updated" {
		t.Fatalf("expected 'updated', got %q", updated.Content)
	}
}

func TestPostRepository_Delete(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user
	userID := utils.NewID()
	user := models.User{ID: userID, Username: "deleteuser", Email: "delete@test.com"}
	db.Create(&user)

	// Create post
	postID := utils.NewID()
	post := models.Post{ID: postID, AuthorID: userID, Content: "to be deleted"}
	repo.Create(&post)

	// Delete
	err := repo.Delete(postID)
	if err != nil {
		t.Fatalf("Delete: %v", err)
	}

	// Verify deleted
	_, err = repo.GetByID(postID)
	if err == nil {
		t.Fatal("expected error for deleted post")
	}
}

func TestPostRepository_SetAndGetReaction(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "reactuser", Email: "react@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	db.Create(&user)
	db.Create(&post)

	// Set reaction
	err := repo.SetPostReaction(userID, postID, 1)
	if err != nil {
		t.Fatalf("SetPostReaction: %v", err)
	}

	// Get reaction
	val, err := repo.GetPostReaction(userID, postID)
	if err != nil {
		t.Fatalf("GetPostReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}

	// Clear reaction
	err = repo.SetPostReaction(userID, postID, 0)
	if err != nil {
		t.Fatalf("Clear reaction: %v", err)
	}

	val, err = repo.GetPostReaction(userID, postID)
	if err != nil {
		t.Fatalf("GetPostReaction after clear: %v", err)
	}
	if val != 0 {
		t.Fatalf("expected 0 after clear, got %d", val)
	}
}

func TestPostRepository_SetAndGetReplyReaction(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user, post, and reply
	userID := utils.NewID()
	postID := utils.NewID()
	replyID := utils.NewID()
	user := models.User{ID: userID, Username: "replyreactuser", Email: "replyreact@test.com"}
	post := models.Post{ID: postID, AuthorID: utils.NewID(), Content: "post"}
	reply := models.Reply{ID: replyID, PostID: postID, AuthorID: utils.NewID(), Content: "reply"}
	db.Create(&user)
	db.Create(&post)
	repo.CreateComment(&reply)

	// Set reaction
	err := repo.SetReplyReaction(userID, replyID, 1)
	if err != nil {
		t.Fatalf("SetReplyReaction: %v", err)
	}

	// Get reaction
	val, err := repo.GetReplyReaction(userID, replyID)
	if err != nil {
		t.Fatalf("GetReplyReaction: %v", err)
	}
	if val != 1 {
		t.Fatalf("expected 1, got %d", val)
	}
}

func TestPostRepository_GetCommentsByPostID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user and post
	userID := utils.NewID()
	postID := utils.NewID()
	user := models.User{ID: userID, Username: "commentuser", Email: "comment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	db.Create(&user)
	db.Create(&post)

	// Create comments
	comment1 := models.Reply{ID: utils.NewID(), PostID: postID, AuthorID: userID, Content: "comment1"}
	comment2 := models.Reply{ID: utils.NewID(), PostID: postID, AuthorID: userID, Content: "comment2"}
	repo.CreateComment(&comment1)
	repo.CreateComment(&comment2)

	// Get comments
	comments, err := repo.GetCommentsByPostID(postID)
	if err != nil {
		t.Fatalf("GetCommentsByPostID: %v", err)
	}
	if len(comments) != 2 {
		t.Fatalf("expected 2 comments, got %d", len(comments))
	}
}

func TestPostRepository_GetCommentByID(t *testing.T) {
	_, db := SetupTestEnv()
	repo := repositories.NewPostRepository(db)

	// Create test user, post, and comment
	userID := utils.NewID()
	postID := utils.NewID()
	commentID := utils.NewID()
	user := models.User{ID: userID, Username: "getcommentuser", Email: "getcomment@test.com"}
	post := models.Post{ID: postID, AuthorID: userID, Content: "post"}
	comment := models.Reply{ID: commentID, PostID: postID, AuthorID: userID, Content: "comment"}
	db.Create(&user)
	db.Create(&post)
	repo.CreateComment(&comment)

	// Get comment
	got, err := repo.GetCommentByID(commentID)
	if err != nil {
		t.Fatalf("GetCommentByID: %v", err)
	}
	if got.Content != "comment" {
		t.Fatalf("expected 'comment', got %q", got.Content)
	}
}
