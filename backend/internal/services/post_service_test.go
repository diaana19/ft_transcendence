package services

import (
	"strings"
	"testing"
	"time"

	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// mockPostRepo is a stateful in-memory PostRepository used only for the
// validation branches below that cannot be reached through the HTTP layer
// (request binding rejects empty/over-long content before the service runs).
// All other post/comment behaviour is covered end-to-end via the /posts API.
type mockPostRepo struct {
	posts         map[string]*models.Post
	comments      map[string]*models.Reply
	postReaction  map[string]int
	replyReaction map[string]int
}

func newMockPostRepo() *mockPostRepo {
	return &mockPostRepo{
		posts:         map[string]*models.Post{},
		comments:      map[string]*models.Reply{},
		postReaction:  map[string]int{},
		replyReaction: map[string]int{},
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

func (m *mockPostRepo) GetByTag(tag string, limit, offset int) ([]models.Post, int64, error) {
	out := make([]models.Post, 0, len(m.posts))
	for _, p := range m.posts {
		for _, t := range p.Tags {
			if t == tag {
				out = append(out, *p)
				break
			}
		}
	}
	return out, int64(len(out)), nil
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

func (m *mockPostRepo) SetPostReaction(userID, postID string, value int) error {
	if value == 0 {
		delete(m.postReaction, userID+"|"+postID)
		return nil
	}
	m.postReaction[userID+"|"+postID] = value
	return nil
}

func (m *mockPostRepo) GetPostReaction(userID, postID string) (int, error) {
	return m.postReaction[userID+"|"+postID], nil
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

func (m *mockPostRepo) SetReplyReaction(userID, replyID string, value int) error {
	if value == 0 {
		delete(m.replyReaction, userID+"|"+replyID)
		return nil
	}
	m.replyReaction[userID+"|"+replyID] = value
	return nil
}

func (m *mockPostRepo) GetReplyReaction(userID, replyID string) (int, error) {
	return m.replyReaction[userID+"|"+replyID], nil
}

func (m *mockPostRepo) TopTags(_ time.Time, _ int) ([]models.TagCount, error) {
	return nil, nil
}

func TestPostService_ContentValidationBranches(t *testing.T) {
	s := NewPostService(newMockPostRepo())

	if _, err := s.CreatePost("", "a1", nil, nil); err == nil {
		t.Fatal("empty post content should error")
	}
	if _, err := s.CreateComment("", "u1", "p1", nil); err == nil {
		t.Fatal("empty comment content should error")
	}
	if _, err := s.CreateComment(strings.Repeat("x", 281), "u1", "p1", nil); err == nil {
		t.Fatal("comment over 280 chars should error")
	}
}
