package services

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// GDPRService handles user data export and deletion for GDPR.
type GDPRService struct {
	db *gorm.DB
}

// NewGDPRService creates a new GDPRService.
func NewGDPRService(db *gorm.DB) *GDPRService {
	return &GDPRService{db: db}
}

// GDPRExportData holds the user data returned by an export.
type GDPRExportData struct {
	User  models.UserResponse   `json:"user"`
	Posts []models.PostResponse `json:"posts"`
}

// ExportUserData returns the user and all the posts of the user.
func (s *GDPRService) ExportUserData(userID string) (*GDPRExportData, error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return nil, err
	}

	var posts []models.Post
	s.db.Where("author_id = ?", userID).Find(&posts)

	postResponses := make([]models.PostResponse, len(posts))
	for i, p := range posts {
		postResponses[i] = p.ToResponse()
	}

	return &GDPRExportData{
		User:  user.ToResponse(),
		Posts: postResponses,
	}, nil
}

// GetUserContact returns the email and username of the user.
func (s *GDPRService) GetUserContact(userID string) (email, username string, err error) {
	var user models.User
	if err := s.db.First(&user, "id = ?", userID).Error; err != nil {
		return "", "", err
	}
	return user.Email, user.Username, nil
}

// DeleteUserData removes the posts and the user account.
func (s *GDPRService) DeleteUserData(userID string) error {
	if err := s.db.Where("author_id = ?", userID).Delete(&models.Post{}).Error; err != nil {
		return err
	}

	return s.db.Where("id = ?", userID).Delete(&models.User{}).Error
}
