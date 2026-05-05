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

func (s *PostService) GetPostsByAuthor(authorID string) ([]models.Post, error) {
	return s.repo.GetByAuthorID(authorID)
}

func (s *PostService) CreatePost(content, authorID string, media *string) (*models.Post, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 280 {
		return nil, errors.New("content must not exceed 280 characters")
	}

	post := &models.Post{
		ID:       uuid.New().String(),
		Content:  content,
		MediaURL: media,
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

func (s *PostService) ToggleLike(userID, postID string) (bool, *models.Post, error) {

	post, err := s.repo.GetByID(postID)
	if err != nil {
		return false, nil, err
	}

	alreadyLiked, err := s.repo.HasLiked(userID, postID)
	if err != nil {
		return false, nil, err
	}

	if alreadyLiked {
		if err := s.repo.UnlikePost(userID, postID); err != nil {
			return false, nil, err
		}
		post.LikesCount--
		if post.LikesCount < 0 {
			post.LikesCount = 0
		}
		return false, post, nil
	}

	if err := s.repo.LikePost(userID, postID); err != nil {
		return false, nil, err
	}
	post.LikesCount++
	return true, post, nil
}


func (s *PostService) HasLiked(userID, postID string) (bool, error) {
	return s.repo.HasLiked(userID, postID)
}



func (s *PostService) CreateComment(content, authorID, postID string) (*models.Reply, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 280 {
		return nil, errors.New("content must not exceed 280 characters")
	}


	if _, err := s.repo.GetByID(postID); err != nil {
		return nil, errors.New("post not found")
	}

	comment := &models.Reply{
		ID:       uuid.New().String(),
		PostID:   postID,
		AuthorID: authorID,
		Content:  content,
	}

	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}


	return s.repo.GetCommentByID(comment.ID)
}

func (s *PostService) GetComments(postID string) ([]models.Reply, error) {

	if _, err := s.repo.GetByID(postID); err != nil {
		return nil, errors.New("post not found")
	}
	return s.repo.GetCommentsByPostID(postID)
}

func (s *PostService) UpdateComment(commentID string, input models.UpdateCommentInput, authorID string) (*models.Reply, error) {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return nil, err
	}
	if comment.AuthorID != authorID {
		return nil, errors.New("you can only update your own comments")
	}
	return s.repo.UpdateComment(commentID, input)
}

func (s *PostService) DeleteComment(commentID, authorID string) error {
	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return err
	}
	if comment.AuthorID != authorID {
		return errors.New("you can only delete your own comments")
	}
	return s.repo.DeleteComment(commentID)
}
