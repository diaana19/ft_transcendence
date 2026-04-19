package services

import (
	"errors"
	"github.com/Transcendence/models"
	"github.com/Transcendence/repositories"
	"github.com/google/uuid"
)

type PostService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPosts(page, limit int) ([]models.Post, int64, error) {
	return s.repo.GetAll(page, limit)
}

func (s *PostService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) CreatePost(content, authorID string) (*models.Post, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 280 {
		return nil, errors.New("content must not exceed 280 characters")
	}

	post := &models.Post{
		ID:       uuid.New().String(),
		Content:  content,
		AuthorID: authorID,
	}

	err := s.repo.Create(post)
	return post, err
}

func (s *PostService) UpdatePost(id string, input models.UpdatePostInput, authorID string) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if post.AuthorID != authorID {
		return nil, errors.New("you can only update your own posts")
	}

	return s.repo.Update(id, input)
}

func (s *PostService) DeletePost(id string, authorID string) error {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if post.AuthorID != authorID {
		return errors.New("you can only delete your own posts")
	}

	return s.repo.Delete(id)
}
