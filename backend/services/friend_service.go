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
	err := s.DB.Where("user_id = ? AND friend_id = ? AND status IN (?)", userID, targetID, []string{"pending", "accepted"}).First(&existing).Error
	if err == nil {
		return errors.New("relationship already exists")
	}
	friend := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   "pending",
	}
	return s.DB.Create(&friend).Error
}

func (s *FriendService) AcceptRequest(userID, requesterID string) error {
	if userID == requesterID {
		return errors.New("cannot accept yourself")
	}

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
	err := s.DB.Where("user_id = ? AND friend_id = ? AND status = ?", userID, targetID, "follow").First(&existing).Error
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

func (s *FriendService) CountFollowers(userID string) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).Where("friend_id = ? AND status = ?", userID, "follow").Count(&count).Error
	return count, err
}

func (s *FriendService) CountFollowing(userID string) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).Where("user_id = ? AND status = ?", userID, "follow").Count(&count).Error
	return count, err
}

func (s *FriendService) Unfollow(userID, targetID string) error {
	result := s.DB.Where("user_id = ? AND friend_id = ? AND status = ?", userID, targetID, "follow").Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not following this user")
	}
	return nil
}

func (s *FriendService) RemoveFriend(userID, targetID string) error {
	result := s.DB.Where("((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?", userID, targetID, targetID, userID, "accepted").Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("you are not friend with this user")
	}
	return nil
}

func (s *FriendService) RejectRequest(userID, targetID string) error {
	result := s.DB.Where("user_id = ? AND friend_id = ? AND status = ?", userID, targetID, "pending").Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no pending request found")
	}
	return nil
}

func (s *FriendService) GetFollowers(userID string) ([]models.User, error) {
	var followers []models.User
	err := s.DB.Joins("JOIN friends ON friends.friend_id = user.id").Where("friends.friend_id = ? AND friends.status = ?", userID, "follow").Find(&followers).Error
	return followers, err
}

func (s *FriendService) GetFollowing(userID string) ([]models.User, error) {
	var following []models.User
	err := s.DB.Joins("JOIN friends ON friends.friend_id = users.id").Where("friends.user_id = ? AND friends.status = ?", userID, "follow").Find(&following).Error
	return following, err
}

func (s *FriendService) GetFriends(userID string) ([]models.User, error) {
	var friends []models.User
	err := s.DB.Joins("JOIN friends ON (friends.user_id = users.id AND friends.friend_id = ?) OR (friends.friend_id = users.id AND friends.user_id = ?)", userID, userID).Where("friends.status = ?", "accepted").Find(&friends).Error
	return friends, err
}
