package services

import (
	"strings"
	"testing"

	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

// mockPostRepo is a stateful in-memory implementation of PostRepository,
// letting the post service be exercised without a database.
type mockPostRepo struct {
	posts    map[string]*models.Post
	comments map[string]*models.Reply
	liked    map[string]bool
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{
		posts:    map[string]*models.Post{},
		comments: map[string]*models.Reply{},
		liked:    map[string]bool{},
	}
}

func (m *mockPostRepo) GetAll(page, limit int) ([]models.Post, int64, error) {
	out := []models.Post{}
	for _, p := range m.posts {
		out = append(out, *p)
	}
	return out, int64(len(out)), nil
}

func (m *mockPostRepo) GetByID(id string) (*models.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return p, nil
}

func (m *mockPostRepo) GetByAuthorID(authorID string) ([]models.Post, error) {
	out := []models.Post{}
	for _, p := range m.posts {
		if p.AuthorID == authorID {
			out = append(out, *p)
		}
	}
	return out, nil
}

func (m *mockPostRepo) Create(post *models.Post) error {
	m.posts[post.ID] = post
	return nil
}

func (m *mockPostRepo) Update(id string, input models.UpdatePostInput) (*models.Post, error) {
	p, ok := m.posts[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	p.Content = input.Content
	return p, nil
}

func (m *mockPostRepo) Delete(id string) error {
	if _, ok := m.posts[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.posts, id)
	return nil
}

func (m *mockPostRepo) LikePost(userID, postID string) error {
	m.liked[userID+"|"+postID] = true
	return nil
}

func (m *mockPostRepo) UnlikePost(userID, postID string) error {
	delete(m.liked, userID+"|"+postID)
	return nil
}

func (m *mockPostRepo) HasLiked(userID, postID string) (bool, error) {
	return m.liked[userID+"|"+postID], nil
}

func (m *mockPostRepo) CreateComment(comment *models.Reply) error {
	m.comments[comment.ID] = comment
	return nil
}

func (m *mockPostRepo) GetCommentsByPostID(postID string) ([]models.Reply, error) {
	out := []models.Reply{}
	for _, c := range m.comments {
		if c.PostID == postID {
			out = append(out, *c)
		}
	}
	return out, nil
}

func (m *mockPostRepo) GetCommentByID(id string) (*models.Reply, error) {
	c, ok := m.comments[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	return c, nil
}

func (m *mockPostRepo) UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error) {
	c, ok := m.comments[id]
	if !ok {
		return nil, gorm.ErrRecordNotFound
	}
	c.Content = input.Content
	return c, nil
}

func (m *mockPostRepo) DeleteComment(id string) error {
	if _, ok := m.comments[id]; !ok {
		return gorm.ErrRecordNotFound
	}
	delete(m.comments, id)
	return nil
}

func TestPostService_CreatePost(t *testing.T) {
	s := NewPostService(newMockPostRepo())

	if _, err := s.CreatePost("", "a1", nil); err == nil {
		t.Fatal("empty content should error")
	}
	if _, err := s.CreatePost(strings.Repeat("x", 281), "a1", nil); err == nil {
		t.Fatal("content over 280 should error")
	}
	p, err := s.CreatePost("hello", "a1", nil)
	if err != nil || p.Content != "hello" || p.AuthorID != "a1" || p.ID == "" {
		t.Fatalf("create failed: %v %+v", err, p)
	}
}

func TestPostService_GetPostsAndGetPost(t *testing.T) {
	repo := newMockPostRepo()
	s := NewPostService(repo)
	repo.posts["p1"] = &models.Post{ID: "p1", AuthorID: "a1", Content: "c"}

	posts, total, err := s.GetPosts(1, 10)
	if err != nil || total != 1 || len(posts) != 1 {
		t.Fatalf("GetPosts: %v total=%d", err, total)
	}
	if _, err := s.GetPost("p1"); err != nil {
		t.Fatalf("GetPost existing: %v", err)
	}
	if _, err := s.GetPost("missing"); err == nil {
		t.Fatal("GetPost missing should error")
	}
	by, err := s.GetPostsByAuthor("a1")
	if err != nil || len(by) != 1 {
		t.Fatalf("GetPostsByAuthor: %v len=%d", err, len(by))
	}
}

func TestPostService_UpdateAndDeleteOwnership(t *testing.T) {
	repo := newMockPostRepo()
	s := NewPostService(repo)
	repo.posts["p1"] = &models.Post{ID: "p1", AuthorID: "owner", Content: "old"}

	if _, err := s.UpdatePost("missing", models.UpdatePostInput{Content: "x"}, "owner"); err == nil {
		t.Fatal("update missing should error")
	}
	if _, err := s.UpdatePost("p1", models.UpdatePostInput{Content: "x"}, "intruder"); err == nil {
		t.Fatal("update by non-owner should error")
	}
	if _, err := s.UpdatePost("p1", models.UpdatePostInput{Content: "new"}, "owner"); err != nil {
		t.Fatalf("owner update: %v", err)
	}

	if err := s.DeletePost("missing", "owner"); err == nil {
		t.Fatal("delete missing should error")
	}
	if err := s.DeletePost("p1", "intruder"); err == nil {
		t.Fatal("delete by non-owner should error")
	}
	if err := s.DeletePost("p1", "owner"); err != nil {
		t.Fatalf("owner delete: %v", err)
	}
}

func TestPostService_ToggleLike(t *testing.T) {
	repo := newMockPostRepo()
	s := NewPostService(repo)
	repo.posts["p1"] = &models.Post{ID: "p1", AuthorID: "a1", LikesCount: 0}

	if _, _, err := s.ToggleLike("u1", "missing"); err == nil {
		t.Fatal("like missing post should error")
	}

	liked, post, err := s.ToggleLike("u1", "p1")
	if err != nil || !liked || post.LikesCount != 1 {
		t.Fatalf("first toggle: liked=%v count=%d err=%v", liked, post.LikesCount, err)
	}
	liked, post, err = s.ToggleLike("u1", "p1")
	if err != nil || liked || post.LikesCount != 0 {
		t.Fatalf("second toggle: liked=%v count=%d err=%v", liked, post.LikesCount, err)
	}

	has, err := s.HasLiked("u1", "p1")
	if err != nil || has {
		t.Fatalf("HasLiked after unlike: has=%v err=%v", has, err)
	}
}

func TestPostService_Comments(t *testing.T) {
	repo := newMockPostRepo()
	s := NewPostService(repo)
	repo.posts["p1"] = &models.Post{ID: "p1", AuthorID: "a1"}

	if _, err := s.CreateComment("", "u1", "p1"); err == nil {
		t.Fatal("empty comment should error")
	}
	if _, err := s.CreateComment(strings.Repeat("x", 281), "u1", "p1"); err == nil {
		t.Fatal("comment over 280 should error")
	}
	if _, err := s.CreateComment("hi", "u1", "missing"); err == nil {
		t.Fatal("comment on missing post should error")
	}
	c, err := s.CreateComment("nice", "u1", "p1")
	if err != nil || c.Content != "nice" {
		t.Fatalf("create comment: %v %+v", err, c)
	}

	comments, err := s.GetComments("p1")
	if err != nil || len(comments) != 1 {
		t.Fatalf("get comments: %v len=%d", err, len(comments))
	}
	if _, err := s.GetComments("missing"); err == nil {
		t.Fatal("get comments missing post should error")
	}

	if _, err := s.UpdateComment(c.ID, models.UpdateCommentInput{Content: "edited"}, "intruder"); err == nil {
		t.Fatal("update comment by non-owner should error")
	}
	if _, err := s.UpdateComment("missing", models.UpdateCommentInput{Content: "x"}, "u1"); err == nil {
		t.Fatal("update missing comment should error")
	}
	if _, err := s.UpdateComment(c.ID, models.UpdateCommentInput{Content: "edited"}, "u1"); err != nil {
		t.Fatalf("owner update comment: %v", err)
	}

	if err := s.DeleteComment("missing", "u1"); err == nil {
		t.Fatal("delete missing comment should error")
	}
	if err := s.DeleteComment(c.ID, "intruder"); err == nil {
		t.Fatal("delete comment by non-owner should error")
	}
	if err := s.DeleteComment(c.ID, "u1"); err != nil {
		t.Fatalf("owner delete comment: %v", err)
	}
}
