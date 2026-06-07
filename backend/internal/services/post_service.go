package services

import (
	"errors"

	"github.com/google/uuid"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
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

// resolveReaction maps a pressed button (+1 like, -1 dislike) against the
// caller's current reaction into the resulting value: pressing the button you
// already have clears it (toggle off), otherwise it becomes your new reaction.
func resolveReaction(current, pressed int) int {
	if current == pressed {
		return 0
	}
	return pressed
}

// ReactToPost applies a like/dislike press to a post and returns the caller's
// resulting reaction (+1/-1/0) along with the post carrying fresh counts.
func (s *PostService) ReactToPost(userID, postID string, pressed int) (int, *models.Post, error) {
	if _, err := s.repo.GetByID(postID); err != nil {
		return 0, nil, err
	}

	current, err := s.repo.GetPostReaction(userID, postID)
	if err != nil {
		return 0, nil, err
	}

	value := resolveReaction(current, pressed)
	if err := s.repo.SetPostReaction(userID, postID, value); err != nil {
		return 0, nil, err
	}

	post, err := s.repo.GetByID(postID)
	if err != nil {
		return 0, nil, err
	}
	return value, post, nil
}

func (s *PostService) GetPostReaction(userID, postID string) (int, error) {
	return s.repo.GetPostReaction(userID, postID)
}

// ReactToComment is the reply counterpart of ReactToPost.
func (s *PostService) ReactToComment(userID, commentID string, pressed int) (int, *models.Reply, error) {
	if _, err := s.repo.GetCommentByID(commentID); err != nil {
		return 0, nil, err
	}

	current, err := s.repo.GetReplyReaction(userID, commentID)
	if err != nil {
		return 0, nil, err
	}

	value := resolveReaction(current, pressed)
	if err := s.repo.SetReplyReaction(userID, commentID, value); err != nil {
		return 0, nil, err
	}

	comment, err := s.repo.GetCommentByID(commentID)
	if err != nil {
		return 0, nil, err
	}
	return value, comment, nil
}

func (s *PostService) GetCommentReaction(userID, commentID string) (int, error) {
	return s.repo.GetReplyReaction(userID, commentID)
}

func (s *PostService) CreateComment(content, authorID, postID string, fileID *string) (*models.Reply, error) {
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
		FileID:   fileID,
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

func (s *PostService) UpdateComment(
	commentID string, input models.UpdateCommentInput, authorID string,
) (*models.Reply, error) {
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
