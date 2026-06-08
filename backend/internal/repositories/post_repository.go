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

func generateUUID() string { return utils.NewID() }

type PostRepository interface {
	GetAll(limit, offset int) ([]models.Post, int64, error)
	GetByID(id string) (*models.Post, error)
	GetByAuthorID(authorID string) ([]models.Post, error)
	GetByTag(tag string, limit, offset int) ([]models.Post, int64, error)
	Create(post *models.Post) error
	Update(id string, input models.UpdatePostInput) (*models.Post, error)
	Delete(id string) error

	TopTags(since time.Time, limit int) ([]models.TagCount, error)

	SetPostReaction(userID, postID string, value int) error
	GetPostReaction(userID, postID string) (int, error)

	CreateComment(comment *models.Reply) error
	GetCommentsByPostID(postID string) ([]models.Reply, error)
	GetCommentByID(id string) (*models.Reply, error)
	UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error)
	DeleteComment(id string) error

	SetReplyReaction(userID, replyID string, value int) error
	GetReplyReaction(userID, replyID string) (int, error)
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetAll(limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	result := r.db.Preload("Author").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)
	r.db.Model(&models.Post{}).Count(&total)

	return posts, total, result.Error
}

func (r *postRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post
	result := r.db.Preload("Author").First(&post, "id = ?", id)
	return &post, result.Error
}

func (r *postRepository) GetByAuthorID(authorID string) ([]models.Post, error) {
	var posts []models.Post
	result := r.db.Preload("Author").Where("author_id = ?", authorID).Order("created_at DESC").Find(&posts)
	return posts, result.Error
}

func (r *postRepository) GetByTag(tag string, limit, offset int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	filter := pq.StringArray{tag}
	r.db.Model(&models.Post{}).Where("tags @> ?", filter).Count(&total)
	result := r.db.Preload("Author").Where("tags @> ?", filter).
		Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts)

	return posts, total, result.Error
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

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

// SetPostReaction reconciles a user's reaction on a post to value (+1 like, -1
// dislike, 0 none), upserting/deleting the single reaction row and adjusting
// the denormalized likes_count / dislikes_count by the transition delta. The
// whole reconciliation runs in one transaction so a row and its counters can
// never drift apart.
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

func (r *postRepository) GetPostReaction(userID, postID string) (int, error) {
	var reaction models.PostReaction
	err := r.db.Where("user_id = ? AND post_id = ?", userID, postID).First(&reaction).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return reaction.Value, err
}

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

func (r *postRepository) GetCommentsByPostID(postID string) ([]models.Reply, error) {
	var comments []models.Reply
	result := r.db.Preload("Author").Preload("File").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments)
	return comments, result.Error
}

func (r *postRepository) GetCommentByID(id string) (*models.Reply, error) {
	var comment models.Reply
	result := r.db.Preload("Author").Preload("File").First(&comment, "id = ?", id)
	return &comment, result.Error
}

func (r *postRepository) UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error) {
	var comment models.Reply
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&comment).Update("content", input.Content).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Author").Preload("File").First(&comment, "id = ?", id)
	return &comment, nil
}

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

// SetReplyReaction is the reply counterpart of SetPostReaction: it reconciles a
// user's reaction on a reply and adjusts that reply's denormalized counters in
// one transaction.
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

func (r *postRepository) GetReplyReaction(userID, replyID string) (int, error) {
	var reaction models.ReplyReaction
	err := r.db.Where("user_id = ? AND reply_id = ?", userID, replyID).First(&reaction).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return 0, nil
	}
	return reaction.Value, err
}
