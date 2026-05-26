/*
** File: gdpr_service.go
** Description: Handles GDPR compliance operations
** Responsibilities:
** - Export all user data in JSON format
** - Permanently delete all user data from the database
*/

package services

import (
	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type GDPRService struct {
	db *gorm.DB
}

func NewGDPRService(db *gorm.DB) *GDPRService {
	return &GDPRService{db: db}
}

type GDPRExportData struct {
	User  models.UserResponse  `json:"user"`
	Posts []models.PostResponse `json:"posts"`
}

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

func (s *GDPRService) DeleteUserData(userID string) error {
	if err := s.db.Where("author_id = ?", userID).Delete(&models.Post{}).Error; err != nil {
		return err
	}

	if err := s.db.Where("id = ?", userID).Delete(&models.User{}).Error; err != nil {
		return err
	}

	return nil
}