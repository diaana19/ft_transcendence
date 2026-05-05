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
	var target models.User
	if err := s.DB.First(&target, "id = ?", targetID).Error; err != nil {
		return errors.New("target user not found")
	}
	var existing models.Friend
	err := s.DB.Where("user_id = ? AND friend_id = ?", userID, targetID).First(&existing).Error
	if err == nil {
		return errors.New("relationship already exists")
	}
	friend := models.Friend{
		UserID: userID,
		FriendID: targetID,
		Status: "pending",
	}
	return s.DB.Create(&friend).Error
}

func (s *FriendService) AcceptRequest(userID, requesterID string) error {
	if userID == requesterID {
		return errors.New("cannot accept yourself")
	}

	// Find the pending request from requester → user
	var friend models.Friend
	err := s.DB.Where(
		"user_id = ? AND friend_id = ? AND status = ?",
		requesterID, userID, "pending",
	).First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("no pending request found")
		}
		return err
	}

	friend.Status = "accepted"
	return s.DB.Save(&friend).Error
}

func (s *FriendService) Follow(userID, targetID string) error {
	if userID == targetID {
		return errors.New("cannot add yourself")
	}
	var target models.User
	if err := s.DB.First(&target, "id = ?", targetID).Error; err != nil {
		return errors.New("target user not found")
	}
	var existing models.Friend
	err := s.DB.Where("user_id = ? AND friend_id = ?", userID, targetID).First(&existing).Error
	if err == nil {
		return errors.New("relationship already exists")
	}
	follow := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   "follow",
	}

	return s.DB.Create(&follow).Error
}
