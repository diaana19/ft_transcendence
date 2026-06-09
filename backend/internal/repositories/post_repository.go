package repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/lib/pq"
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/utils"
)

// generateUUID returns a new unique id.
func generateUUID() string { return utils.NewID() }

// PostRepository handles the posts, comments and reactions in the database.
type PostRepository interface {
	// GetAll returns a page of posts and the total count.
	GetAll(limit, offset int) ([]models.Post, int64, error)
	// GetByID returns the post with this id.
	GetByID(id string) (*models.Post, error)
	// GetByAuthorID returns all the posts of an author.
	GetByAuthorID(authorID string) ([]models.Post, error)
	// GetByTag returns a page of posts with this tag and the total count.
	GetByTag(tag string, limit, offset int) ([]models.Post, int64, error)
	// GetRepliedByUser returns a page of posts the user replied to and the total count.
	GetRepliedByUser(userID string, limit, offset int) ([]models.Post, int64, error)
	// Create saves a new post.
	Create(post *models.Post) error
	// Update changes the post fields from the input and returns the updated post.
	Update(id string, input models.UpdatePostInput) (*models.Post, error)
	// Delete removes the post with this id.
	Delete(id string) error

	// TopTags returns the most used tags since the given time.
	TopTags(since time.Time, limit int) ([]models.TagCount, error)

	// SetPostReaction sets the reaction of the user on a post.
	SetPostReaction(userID, postID string, value int) error
	// GetPostReaction returns the reaction of the user on a post.
	GetPostReaction(userID, postID string) (int, error)

	// CreateComment saves a new comment on a post.
	CreateComment(comment *models.Reply) error
	// GetCommentsByPostID returns the comments of a post.
	GetCommentsByPostID(postID string) ([]models.Reply, error)
	// GetCommentByID returns the comment with this id.
	GetCommentByID(id string) (*models.Reply, error)
	// UpdateComment changes the comment fields from the input and returns the updated comment.
	UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error)
	// DeleteComment removes the comment with this id.
	DeleteComment(id string) error

	// SetReplyReaction sets the reaction of the user on a comment.
	SetReplyReaction(userID, replyID string, value int) error
	// GetReplyReaction returns the reaction of the user on a comment.
	GetReplyReaction(userID, replyID string) (int, error)
}

type postRepository struct {
	db *gorm.DB
}

// NewPostRepository creates a new PostRepository using the given database.
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// GetAll returns a page of posts and the total count.
func (r *postRepository) GetAll(limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	result := r.db.Preload("Author").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)
	r.db.Model(&models.Post{}).Count(&total)

	return posts, total, result.Error
}

// GetByID returns the post with this id.
func (r *postRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post
	result := r.db.Preload("Author").First(&post, "id = ?", id)
	return &post, result.Error
}

// GetByAuthorID returns all the posts of an author.
func (r *postRepository) GetByAuthorID(authorID string) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("Author").Where("author_id = ?", authorID).Order("created_at DESC").Find(&posts)
	return posts, result.Error
}

// GetRepliedByUser returns a page of posts the user replied to and the total count.
func (r *postRepository) GetRepliedByUser(userID string, limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	repliedPostIDs := r.db.Model(&models.Reply{}).Select("post_id").Where("author_id = ?", userID)

	r.db.Model(&models.Post{}).Where("id IN (?)", repliedPostIDs).Count(&total)
	result := r.db.Preload("Author").Where("id IN (?)", repliedPostIDs).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)

	return posts, total, result.Error
}

// GetByTag returns a page of posts with this tag and the total count.
func (r *postRepository) GetByTag(tag string, limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	filter := pq.StringArray{tag}
	r.db.Model(&models.Post{}).Where("tags @> ?", filter).Count(&total)
	result := r.db.Preload("Author").Where("tags @> ?", filter).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)

	return posts, total, result.Error
}

// Create saves a new post.
func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

// Update changes the post fields from the input and returns the updated post.
func (r *postRepository) Update(id string, input models.UpdatePostInput) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&post).Select("content", "media_url", "media_mime", "tags").Updates(map[string]any{
		"content":    input.Content,
		"media_url":  input.MediaURL,
		"media_mime": input.MediaMIME,
		"tags":       pq.StringArray(utils.ExtractHashtags(input.Content)),
	}).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Author").First(&post, "id = ?", id)
	return &post, nil
}

// TopTags returns the most used tags since the given time.
func (r *postRepository) TopTags(since time.Time, limit int) ([]models.TagCount, error) {
	var out []models.TagCount
	err := r.db.Raw(`
		SELECT tag, COUNT(*) AS count
		FROM posts, unnest(tags) AS tag
		WHERE deleted_at IS NULL AND created_at >= ?
		GROUP BY tag
		ORDER BY count DESC, tag ASC
		LIMIT ?`, since, limit).Scan(&out).Error
	return out, err
}

// Delete removes the post with this id.
func (r *postRepository) Delete(id string) error {
	result := r.db.Delete(&models.Post{}, "id = ?", id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// reactionDelta returns +1 when the bucket starts to be used, -1 when it stops, else 0.
func reactionDelta(oldValue, newValue, bucket int) int {
	d := 0
	if newValue == bucket {
		d++
	}
	if oldValue == bucket {
		d--
	}
	return d
}

// SetPostReaction sets the reaction of the user on a post.
//
//nolint:dupl // Intentionally parallel to SetReplyReaction; different model types.
func (r *postRepository) SetPostReaction(userID, postID string, value int) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.PostReaction
		err := tx.Where("user_id = ? AND post_id = ?", userID, postID).First(&existing).Error
		found := err == nil
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		oldValue := 0
		if found {
			oldValue = existing.Value
		}
		if oldValue == value {
			return nil
		}

		switch {
		case value == 0:
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		case found:
			if err := tx.Model(&existing).Update("value", value).Error; err != nil {
				return err
			}
		default:
			if err := tx.Create(&models.PostReaction{
				ID: generateUUID(), UserID: userID, PostID: postID, Value: value,
			}).Error; err != nil {
				return err
			}
		}

		return adjustReactionCounts(tx, &models.Post{}, postID, oldValue, value)
	})
	if err != nil {
		return fmt.Errorf("set post reaction: %w", err)
	}
	return nil
}

// GetPostReaction returns the reaction of the user on a post.
func (r *postRepository) GetPostReaction(userID, postID string) (int, error) {
	var reaction models.PostReaction
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&reaction).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return reaction.Value, err
}

// adjustReactionCounts updates the likes and dislikes counts after a reaction change.
func adjustReactionCounts(tx *gorm.DB, model any, id string, oldValue, newValue int) error {
	if d := reactionDelta(oldValue, newValue, 1); d != 0 {
		if err := tx.Model(model).Where("id = ?", id).
			UpdateColumn("likes_count", gorm.Expr("GREATEST(likes_count + ?, 0)", d)).Error; err != nil {
			return err
		}
	}
	if d := reactionDelta(oldValue, newValue, -1); d != 0 {
		if err := tx.Model(model).Where("id = ?", id).
			UpdateColumn("dislikes_count", gorm.Expr("GREATEST(dislikes_count + ?, 0)", d)).Error; err != nil {
			return err
		}
	}
	return nil
}

// CreateComment saves a new comment on a post.
func (r *postRepository) CreateComment(comment *models.Reply) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&models.Post{}).Where("id = ?", comment.PostID).
			UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error
	})
	if err != nil {
		return fmt.Errorf("create comment: %w", err)
	}
	return nil
}

// GetCommentsByPostID returns the comments of a post.
func (r *postRepository) GetCommentsByPostID(postID string) ([]models.Reply, error) {
	var comments []models.Reply
	result := r.db.Preload("Author").Preload("File").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments)
	return comments, result.Error
}

// GetCommentByID returns the comment with this id.
func (r *postRepository) GetCommentByID(id string) (*models.Reply, error) {
	var comment models.Reply
	result := r.db.Preload("Author").Preload("File").First(&comment, "id = ?", id)
	return &comment, result.Error
}

// UpdateComment changes the comment fields from the input and returns the updated comment.
func (r *postRepository) UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error) {
	var comment models.Reply
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}

	updates := map[string]any{"content": input.Content}
	switch {
	case input.NewFileID != nil:
		updates["file_id"] = *input.NewFileID
	case input.RemoveFile:
		updates["file_id"] = nil
	}

	if err := r.db.Model(&comment).Updates(updates).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Author").Preload("File").First(&comment, "id = ?", id)
	return &comment, nil
}

// DeleteComment removes the comment and its reactions, and updates the comment count.
func (r *postRepository) DeleteComment(id string) error {
	var comment models.Reply
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return err
	}
	err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Reply{}, "id = ?", id).Error; err != nil {
			return err
		}
		if err := tx.Where("reply_id = ?", id).Delete(&models.ReplyReaction{}).Error; err != nil {
			return err
		}
		return tx.Model(&models.Post{}).Where("id = ? AND comments_count > 0", comment.PostID).
			UpdateColumn("comments_count", gorm.Expr("comments_count - 1")).Error
	})
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}
	return nil
}

// SetReplyReaction sets the reaction of the user on a comment.
//
//nolint:dupl // Intentionally parallel to SetPostReaction; different model types.
func (r *postRepository) SetReplyReaction(userID, replyID string, value int) error {
	err := r.db.Transaction(func(tx *gorm.DB) error {
		var existing models.ReplyReaction
		err := tx.Where("user_id = ? AND reply_id = ?", userID, replyID).First(&existing).Error
		found := err == nil
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		oldValue := 0
		if found {
			oldValue = existing.Value
		}
		if oldValue == value {
			return nil
		}

		switch {
		case value == 0:
			if err := tx.Delete(&existing).Error; err != nil {
				return err
			}
		case found:
			if err := tx.Model(&existing).Update("value", value).Error; err != nil {
				return err
			}
		default:
			if err := tx.Create(&models.ReplyReaction{
				ID: generateUUID(), UserID: userID, ReplyID: replyID, Value: value,
			}).Error; err != nil {
				return err
			}
		}

		return adjustReactionCounts(tx, &models.Reply{}, replyID, oldValue, value)
	})
	if err != nil {
		return fmt.Errorf("set reply reaction: %w", err)
	}
	return nil
}

// GetReplyReaction returns the reaction of the user on a comment.
func (r *postRepository) GetReplyReaction(userID, replyID string) (int, error) {
	var reaction models.ReplyReaction
	err := r.db.Where("user_id = ? AND reply_id = ?", userID, replyID).First(&reaction).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return reaction.Value, err
}
