package services

import (
	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
	"ft_transcendence/backend/internal/repositories"
)

// GDPRService handles user data export and deletion for GDPR.
type GDPRService struct {
	db       *gorm.DB
	userRepo repositories.UserRepository
}

// NewGDPRService creates a new GDPRService.
func NewGDPRService(db *gorm.DB, userRepo repositories.UserRepository) *GDPRService {
	return &GDPRService{db: db, userRepo: userRepo}
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

// DeleteUserData removes the user and everything they created via the shared user-repository cascade.
func (s *GDPRService) DeleteUserData(userID string) error {
	return s.userRepo.Delete(userID)
}
