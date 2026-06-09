package repositories

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// FileRepository handles the file records in the database.
type FileRepository interface {
	// Create saves a new file.
	Create(file *models.File) error
	// GetByID returns the file with this id.
	GetByID(id string) (*models.File, error)
	// DeleteByOwner deletes the file only if it belongs to the owner.
	DeleteByOwner(fileID, ownerID string) error

	// GrantAccess gives the user access to the file.
	GrantAccess(fileID, userID string) error
	// HasAccess returns true if the user can access the file.
	HasAccess(fileID, userID string) (bool, error)
}

type fileRepository struct {
	db *gorm.DB
}

// NewFileRepository creates a new FileRepository using the given database.
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db: db}
}

// Create saves a new file.
func (r *fileRepository) Create(file *models.File) error {
	return r.db.Create(file).Error
}

// GetByID returns the file with this id.
func (r *fileRepository) GetByID(id string) (*models.File, error) {
	var file models.File
	err := r.db.First(&file, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

// DeleteByOwner deletes the file only if it belongs to the owner.
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

// GrantAccess gives the user access to the file.
func (r *fileRepository) GrantAccess(fileID, userID string) error {
	access := models.FileAccess{
		FileID: fileID,
		UserID: userID,
	}
	return r.db.Where("file_id = ? AND user_id = ?", fileID, userID).FirstOrCreate(&access).Error
}

// HasAccess returns true if the user can access the file.
func (r *fileRepository) HasAccess(fileID, userID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.FileAccess{}).Where("file_id = ? AND user_id = ?", fileID, userID).Count(&count).Error
	return count > 0, err
}
