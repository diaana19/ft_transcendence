package repositories

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	GetAll(page, limit int) ([]models.Post, int64, error)
	GetByID(id string) (*models.Post, error)
	GetByAuthorID(authorID string) ([]models.Post, error)
	Create(post *models.Post) error
	Update(id string, input models.UpdatePostInput) (*models.Post, error)
	Delete(id string) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) GetAll(page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	offset := (page - 1) * limit
	result := r.db.Preload("Author").Offset(offset).Limit(limit).Find(&posts)

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
	result := r.db.Preload("Author").Where("author_id = ?", authorID).Find(&posts)
	return posts, result.Error
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) Update(id string, input models.UpdatePostInput) (*models.Post, error) {
	var post models.Post
	result := r.db.First(&post, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	result = r.db.Model(&post).Updates(input)
	if result.Error != nil {
		return nil, result.Error
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
