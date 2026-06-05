package services

import (
	"strings"
	"testing"

	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// mockPostRepo is a stateful in-memory PostRepository used only for the
// validation branches below that cannot be reached through the HTTP layer
// (request binding rejects empty/over-long content before the service runs).
// All other post/comment behaviour is covered end-to-end via the /posts API.
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

func (m *mockPostRepo) GetAll(_, _ int) ([]models.Post, int64, error) {
	out := make([]models.Post, 0, len(m.posts))
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
	out := make([]models.Post, 0, len(m.posts))
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

func TestPostService_ContentValidationBranches(t *testing.T) {
	s := NewPostService(newMockPostRepo())

	// CreateCommentInput binds min=1/max=280, so these service guards are only
	// reachable below the controller.
	if _, err := s.CreatePost("", "a1", nil); err == nil {
		t.Fatal("empty post content should error")
	}
	if _, err := s.CreateComment("", "u1", "p1"); err == nil {
		t.Fatal("empty comment content should error")
	}
	if _, err := s.CreateComment(strings.Repeat("x", 281), "u1", "p1"); err == nil {
		t.Fatal("comment over 280 chars should error")
	}
}
