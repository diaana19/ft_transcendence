package services

import (
	"errors"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
	"ft_transcendence/backend/internal/utils"
)

// ErrAuthorNotFound is returned when a post is created for an author that does not exist.
var ErrAuthorNotFound = errors.New("author does not exist")

// PostService handles posts, comments and reactions.
type PostService struct {
	repo repositories.PostRepository
}

// NewPostService creates a new PostService.
func NewPostService(repo repositories.PostRepository) *PostService {
	return &PostService{repo: repo}
}

// GetPosts returns a page of posts and the total count.
func (s *PostService) GetPosts(limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetAll(limit, offset)
}

// GetPost returns one post by ID.
func (s *PostService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

// GetPostsByAuthor returns all posts of the author.
func (s *PostService) GetPostsByAuthor(authorID string) ([]models.Post, error) {
	return s.repo.GetByAuthorID(authorID)
}

// GetPostsByTag returns a page of posts with the tag and the total count.
func (s *PostService) GetPostsByTag(tag string, limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetByTag(tag, limit, offset)
}

// GetRepliedPosts returns a page of posts the user replied to and the total count.
func (s *PostService) GetRepliedPosts(userID string, limit, offset int) ([]models.Post, int64, error) {
	return s.repo.GetRepliedByUser(userID, limit, offset)
}

// CreatePost creates a new post. Content is required and must not exceed 280 characters.
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

// GetTrends returns the top tags from the last 7 days.
func (s *PostService) GetTrends(limit int) ([]models.TagCount, error) {
	return s.repo.TopTags(time.Now().Add(-trendWindow), limit)
}

// SearchTags returns a page of tags matching the query and the total count.
func (s *PostService) SearchTags(query string, limit, offset int) ([]models.TagCount, int64, error) {
	return s.repo.SearchTags(query, limit, offset)
}

// UpdatePost updates a post. Only the author can update its own post.
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

// DeletePost deletes a post. Only the author can delete its own post.
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

// resolveReaction returns 0 when the pressed reaction is the same as the current one, else the pressed one.
func resolveReaction(current, pressed int) int {
	if current == pressed {
		return 0
	}
	return pressed
}

// ReactToPost sets or toggles the reaction of the user on a post. It returns the new value and the post.
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

// GetPostReaction returns the reaction of the user on a post.
func (s *PostService) GetPostReaction(userID, postID string) (int, error) {
	return s.repo.GetPostReaction(userID, postID)
}

// ReactToComment sets or toggles the reaction of the user on a comment. It returns the new value and the comment.
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

// GetCommentReaction returns the reaction of the user on a comment.
func (s *PostService) GetCommentReaction(userID, commentID string) (int, error) {
	return s.repo.GetReplyReaction(userID, commentID)
}

// CreateComment creates a comment on a post. Content is required and must not exceed 280 characters.
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

// GetComments returns the comments of a post.
func (s *PostService) GetComments(postID string) ([]models.Reply, error) {
	if _, err := s.repo.GetByID(postID); err != nil {
		return nil, errors.New("post not found")
	}
	return s.repo.GetCommentsByPostID(postID)
}

// UpdateComment updates a comment. Only the author can update its own comment.
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

// DeleteComment deletes a comment. Only the author can delete its own comment.
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
