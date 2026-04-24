package services

import (
	"errors"

	"github.com/Transcendence/models"
	"gorm.io/gorm"
)

type FriendService struct {
	DB *gorm.DB
}

func (s *FriendService) SendRequest(userID, targetID string) error {
	if userID == targetID {
		return errors.New("cannot add yourself")
	}

	friend := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   "pending",
	}

	return s.DB.Create(&friend).Error
}

func (s *FriendService) AcceptRequest(userID, requesterID string) error {
	var friend models.Friend

	err := s.DB.Where("user_id = ? AND friend_id = ? AND status = ?", requesterID, userID, "pending").
		First(&friend).Error

	if err != nil {
		return err
	}

	friend.Status = "accepted"
	return s.DB.Save(&friend).Error
}

func (s *FriendService) Follow(userID, targetID string) error {
	follow := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   "follow",
	}

	return s.DB.Create(&follow).Error
}
