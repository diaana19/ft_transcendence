package services

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

var ErrAuthorNotFound = errors.New("author does not exist")

type PostService struct {
	repo repositories.PostRepository
}

func NewPostService(repo repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

func (s *PostService) GetPosts(limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetAll(limit, offset)
}

func (s *PostService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *PostService) GetPostsByAuthor(authorID string) ([]models.Post, error) {
	return s.repo.GetByAuthorID(authorID)
}

func (s *PostService) GetPostsByTag(tag string, limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetByTag(tag, limit, offset)
}

func (s *PostService) GetRepliedPosts(userID string, limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetRepliedByUser(userID, limit, offset)
}

func (s *PostService) CreatePost(content, authorID string, media *string, mediaMIME *string) (*models.Post, error) {
	if content == "" {
		return nil, errors.New("content is required")
	}
	if len(content) > 280 {
		return nil, errors.New("content must not exceed 280 characters")
	}

	post := &models.Post{
		ID:        utils.NewID(),
		Content:   content,
		MediaURL:  media,
		MediaMIME: mediaMIME,
		AuthorID:  authorID,
		Tags:      utils.ExtractHashtags(content),
	}

	if err := s.repo.Create(post); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23503" {
			return nil, ErrAuthorNotFound
		}
		return nil, err
	}
	return post, nil
}

const trendWindow = 7 * 24 * time.Hour

func (s *PostService) GetTrends(limit int) ([]models.TagCount, error) {
	return s.repo.TopTags(time.Now().Add(-trendWindow), limit)
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

func resolveReaction(current, pressed int) int {
	if current == pressed {
		return 0
	}
	return pressed
}

func (s *PostService) ReactToPost(userID, postID string, pressed int) (int, *models.Post, error) {
	if _, err := s.repo.GetByID(postID); err != nil {
		return 0, nil, err
	}

	current, err := s.repo.GetPostReaction(userID, postID)
	if err != nil {
		return 0, nil, err
	}

	value := resolveReaction(current, pressed)
	if err = s.repo.SetPostReaction(userID, postID, value); err != nil {
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

func (s *PostService) ReactToComment(userID, commentID string, pressed int) (int, *models.Reply, error) {
	if _, err := s.repo.GetCommentByID(commentID); err != nil {
		return 0, nil, err
	}

	current, err := s.repo.GetReplyReaction(userID, commentID)
	if err != nil {
		return 0, nil, err
	}

	value := resolveReaction(current, pressed)
	if err = s.repo.SetReplyReaction(userID, commentID, value); err != nil {
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
		ID:       utils.NewID(),
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
