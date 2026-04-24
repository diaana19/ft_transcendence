package repositories

import (
	"github.com/Transcendence/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func generateUUID() string { return uuid.New().String() }

type PostRepository interface {
	// Posts
	GetAll(page, limit int) ([]models.Post, int64, error)
	GetByID(id string) (*models.Post, error)
	GetByAuthorID(authorID string) ([]models.Post, error)
	Create(post *models.Post) error
	Update(id string, input models.UpdatePostInput) (*models.Post, error)
	Delete(id string) error

	// Likes
	LikePost(userID, postID string) error
	UnlikePost(userID, postID string) error
	HasLiked(userID, postID string) (bool, error)

	// Comments
	CreateComment(comment *models.Reply) error
	GetCommentsByPostID(postID string) ([]models.Reply, error)
	GetCommentByID(id string) (*models.Reply, error)
	UpdateComment(id string, input models.UpdateCommentInput) (*models.Reply, error)
	DeleteComment(id string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

// ─── Posts ────────────────────────────────────────────────────────────────────

func (r *postRepository) GetAll(page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * limit
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

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(id string, input models.UpdatePostInput) (*models.Post, error) {
	var post models.Post
	if err := r.db.First(&post, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if err := r.db.Model(&post).Updates(input).Error; err != nil {
		return nil, err
	}
	r.db.Preload("Author").First(&post, "id = ?", id)
	return &post, nil
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

// ─── Likes ────────────────────────────────────────────────────────────────────

// LikePost creates a Like record and atomically increments the post counter.
func (r *postRepository) LikePost(userID, postID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		like := models.Like{
			ID:     generateUUID(),
			UserID: userID,
			PostID: postID,
		}
		if err := tx.Create(&like).Error; err != nil {
			return err // unique index violation if already liked
		}
		return tx.Model(&models.Post{}).Where("id = ?", postID).
			UpdateColumn("likes_count", gorm.Expr("likes_count + 1")).Error
	})
}

// UnlikePost deletes the Like record and decrements the post counter.
func (r *postRepository) UnlikePost(userID, postID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		result := tx.Where("user_id = ? AND post_id = ?", userID, postID).Delete(&models.Like{})
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return tx.Model(&models.Post{}).Where("id = ? AND likes_count > 0", postID).
			UpdateColumn("likes_count", gorm.Expr("likes_count - 1")).Error
	})
}

// HasLiked returns true when userID has already liked postID.
func (r *postRepository) HasLiked(userID, postID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).
		Where("user_id = ? AND post_id = ?", userID, postID).
		Count(&count).Error
	return count > 0, err
}

// ─── Comments ─────────────────────────────────────────────────────────────────

func (r *postRepository) CreateComment(comment *models.Reply) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(comment).Error; err != nil {
			return err
		}
		return tx.Model(&models.Post{}).Where("id = ?", comment.PostID).
			UpdateColumn("comments_count", gorm.Expr("comments_count + 1")).Error
	})
}

func (r *postRepository) GetCommentsByPostID(postID string) ([]models.Reply, error) {
	var comments []models.Reply
	result := r.db.Preload("Author").
		Where("post_id = ?", postID).
		Order("created_at ASC").
		Find(&comments)
	return comments, result.Error
}

func (r *postRepository) GetCommentByID(id string) (*models.Reply, error) {
	var comment models.Reply
	result := r.db.Preload("Author").First(&comment, "id = ?", id)
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
	r.db.Preload("Author").First(&comment, "id = ?", id)
	return &comment, nil
}

func (r *postRepository) DeleteComment(id string) error {
	// Fetch comment first to get PostID for counter decrement
	var comment models.Reply
	if err := r.db.First(&comment, "id = ?", id).Error; err != nil {
		return err
	}
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Reply{}, "id = ?", id).Error; err != nil {
			return err
		}
		return tx.Model(&models.Post{}).Where("id = ? AND comments_count > 0", comment.PostID).
			UpdateColumn("comments_count", gorm.Expr("comments_count - 1")).Error
	})
}
