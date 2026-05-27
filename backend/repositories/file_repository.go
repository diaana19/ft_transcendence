package repositories

import (
	"gorm.io/gorm"

	"github.com/Transcendence/models"
)

type FileRepository interface {
	Create(file *models.File) error
	GetByID(id string) (*models.File, error)
	DeleteByOwner(fileID, ownerID string) error

	GrantAccess(fileID, userID string) error
	HasAccess(fileID, userID string) (bool, error)
}

type fileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) GetByID(id string) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (r *fileRepository) DeleteByOwner(fileID, userID string) error {
	result := r.db.Where("id = ? AND owner_id = ?", fileID, userID).Delete(&models.File{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *fileRepository) GrantAccess(fileID, userID string) error {
	access := models.FileAccess{
		FileID: fileID,
		UserID: userID,
	}
	return r.db.Where("file_id = ? AND user_id = ?", fileID, userID).FirstOrCreate(&access).Error
}

func (r *fileRepository) HasAccess(fileID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FileAccess{}).Where("file_id = ? AND user_id = ?", fileID, userID).Count(&count).Error
	return count > 0, err
}
