package services

import (
	"errors"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"ft_transcendence/backend/internal/models"
)

const (
	statusPending  = "pending"
	statusAccepted = "accepted"
	statusFollow   = "follow"

	sqlFriendBidirectional = "((user_id = ? AND friend_id = ?) OR " +
		"(user_id = ? AND friend_id = ?)) AND status = ?"
	sqlFriendBidirectionalStatuses = "((user_id = ? AND friend_id = ?) OR " +
		"(user_id = ? AND friend_id = ?)) AND status IN (?)"
	sqlFriendsMutualJoin = "JOIN friends ON (friends.user_id = users.id AND friends.friend_id = ?) " +
		"OR (friends.friend_id = users.id AND friends.user_id = ?)"
)

// FriendService handles friend requests, follows and friend queries.
type FriendService struct {
	DB *gorm.DB
}

// SendRequest creates a pending friend request from the user to the target.
// If the target already sent a pending request to the user, it auto-accepts instead.
func (s *FriendService) SendRequest(userID, targetID string) error {
	// Adding yourself is a no-op rather than an error; the frontend refreshes.
	if userID == targetID {
		return nil
	}
	var target models.User
	if err := s.DB.First(&target, "id = ?", targetID).Error; err != nil {
		return errors.New("target user not found")
	}
	// If target already sent us a pending request, accept it instead of creating a duplicate.
	var reverse models.Friend
	err := s.DB.
		Where("user_id = ? AND friend_id = ? AND status = ?", targetID, userID, statusPending).
		First(&reverse).Error
	if err == nil {
		reverse.Status = statusAccepted
		return s.DB.Save(&reverse).Error
	}
	// Insert the pending request. The unique index on (user_id, friend_id, status)
	// makes this a no-op if the same relationship already exists, so a repeated
	// click reports success and the frontend just refreshes its state.
	friend := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   statusPending,
	}
	return s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&friend).Error
}

// AcceptRequest accepts a pending request sent by the requester to the user.
func (s *FriendService) AcceptRequest(userID, requesterID string) error {
	if userID == requesterID {
		return errors.New("cannot accept yourself")
	}

	var friend models.Friend
	err := s.DB.Where(
		"user_id = ? AND friend_id = ? AND status = ?",
		requesterID, userID, statusPending,
	).First(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Already accepted (e.g. auto-accepted via SendRequest) — not an error.
			already, _ := s.AreFriends(userID, requesterID)
			if already {
				return nil
			}
			return errors.New("no pending request found")
		}
		return err
	}

	friend.Status = statusAccepted
	return s.DB.Save(&friend).Error
}

// Follow makes the user follow the target.
func (s *FriendService) Follow(userID, targetID string) error {
	// Following yourself is a no-op rather than an error; the frontend refreshes.
	if userID == targetID {
		return nil
	}
	var target models.User
	if err := s.DB.First(&target, "id = ?", targetID).Error; err != nil {
		return errors.New("target user not found")
	}
	// Insert the follow. The unique index on (user_id, friend_id, status) makes
	// this a no-op if it already exists, so a repeated click reports success and
	// the frontend just refreshes its state.
	follow := models.Friend{
		UserID:   userID,
		FriendID: targetID,
		Status:   statusFollow,
	}

	return s.DB.Clauses(clause.OnConflict{DoNothing: true}).Create(&follow).Error
}

// CountFriends returns how many accepted friend relationships the given user has.
func (s *FriendService) CountFriends(userID string) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).
		Where("(user_id = ? OR friend_id = ?) AND status = ?", userID, userID, statusAccepted).
		Count(&count).Error
	return count, err
}

// CountFollowers returns how many users follow the given user.
func (s *FriendService) CountFollowers(userID string) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).Where("friend_id = ? AND status = ?", userID, statusFollow).Count(&count).Error
	return count, err
}

// CountFollowing returns how many users the given user follows.
func (s *FriendService) CountFollowing(userID string) (int64, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).Where("user_id = ? AND status = ?", userID, statusFollow).Count(&count).Error
	return count, err
}

// Unfollow removes the follow relationship only, leaving any friend relationship intact.
func (s *FriendService) Unfollow(userID, targetID string) error {
	result := s.DB.
		Where("user_id = ? AND friend_id = ? AND status = ?", userID, targetID, statusFollow).
		Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("not following this user")
	}
	return nil
}

// RemoveFriend deletes the friend relationship between the two users.
func (s *FriendService) RemoveFriend(userID, targetID string) error {
	result := s.DB.
		Where(sqlFriendBidirectionalStatuses, userID, targetID, targetID, userID, []string{statusPending, statusAccepted}).
		Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("you are not friend with this user")
	}
	return nil
}

// RejectRequest removes a pending request sent by the requester to the user.
func (s *FriendService) RejectRequest(userID, requesterID string) error {
	result := s.DB.Where(
		"user_id = ? AND friend_id = ? AND status = ?",
		requesterID, userID, statusPending,
	).Delete(&models.Friend{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no pending request found")
	}
	return nil
}

// GetFollowers returns the users that follow the given user.
func (s *FriendService) GetFollowers(userID string) ([]models.User, error) {
	var followers []models.User
	err := s.DB.
		Joins("JOIN friends ON friends.user_id = users.id").
		Where("friends.friend_id = ? AND friends.status = ?", userID, statusFollow).
		Find(&followers).Error
	return followers, err
}

// GetFollowing returns the users that the given user follows.
func (s *FriendService) GetFollowing(userID string) ([]models.User, error) {
	var following []models.User
	err := s.DB.
		Joins("JOIN friends ON friends.friend_id = users.id").
		Where("friends.user_id = ? AND friends.status = ?", userID, statusFollow).
		Find(&following).Error
	return following, err
}

// GetFriends returns the accepted friends of the user.
func (s *FriendService) GetFriends(userID string) ([]models.User, error) {
	var friends []models.User
	err := s.DB.
		Joins(sqlFriendsMutualJoin, userID, userID).
		Where("friends.status = ?", statusAccepted).
		Find(&friends).Error
	return friends, err
}

// AreFriends returns true if the two users have an accepted friend relationship.
func (s *FriendService) AreFriends(userID1, userID2 string) (bool, error) {
	var count int64
	err := s.DB.Model(&models.Friend{}).
		Where(sqlFriendBidirectional, userID1, userID2, userID2, userID1, statusAccepted).
		Count(&count).Error
	return count > 0, err
}
