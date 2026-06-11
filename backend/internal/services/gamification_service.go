package services

import (
	"math/bits"

	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// Metric holds a count and the level computed from it.
type Metric struct {
	Level int   `json:"level"`
	Count int64 `json:"count"`
}

// asMetric builds a Metric from a count.
func asMetric(count int64) Metric {
	return Metric{Level: Level(count), Count: count}
}

// GamificationStats holds the level and the metrics of a user.
type GamificationStats struct {
	Level     int    `json:"level"`
	Total     int64  `json:"total"`
	Posts     Metric `json:"posts"`
	Likes     Metric `json:"likes"`
	Friends   Metric `json:"friends"`
	Followers Metric `json:"followers"`
	Following Metric `json:"following"`
}

// GamificationService computes user stats and the leaderboard.
type GamificationService struct {
	db     *gorm.DB
	friend *FriendService
}

// NewGamificationService creates a new GamificationService.
func NewGamificationService(db *gorm.DB, friend *FriendService) *GamificationService {
	return &GamificationService{db: db, friend: friend}
}

// Compute returns the gamification stats of the user from posts, likes, followers and following.
func (s *GamificationService) Compute(userID string) (GamificationStats, error) {
	var posts int64
	if err := s.db.Model(&models.Post{}).
		Where("author_id = ?", userID).
		Count(&posts).Error; err != nil {
		return GamificationStats{}, err
	}

	var postLikes int64
	if err := s.db.Model(&models.Post{}).
		Where("author_id = ?", userID).
		Select("COALESCE(SUM(likes_count), 0)").
		Scan(&postLikes).Error; err != nil {
		return GamificationStats{}, err
	}

	var replyLikes int64
	if err := s.db.Model(&models.Reply{}).
		Where("author_id = ?", userID).
		Select("COALESCE(SUM(likes_count), 0)").
		Scan(&replyLikes).Error; err != nil {
		return GamificationStats{}, err
	}
	likesReceived := postLikes + replyLikes

	friends, err := s.friend.CountFriends(userID)
	if err != nil {
		return GamificationStats{}, err
	}
	followers, err := s.friend.CountFollowers(userID)
	if err != nil {
		return GamificationStats{}, err
	}
	following, err := s.friend.CountFollowing(userID)
	if err != nil {
		return GamificationStats{}, err
	}

	total := posts + likesReceived + friends + followers + following
	return GamificationStats{
		Level:     Level(total),
		Total:     total,
		Posts:     asMetric(posts),
		Likes:     asMetric(likesReceived),
		Friends:   asMetric(friends),
		Followers: asMetric(followers),
		Following: asMetric(following),
	}, nil
}

// Level returns the level for a total. It grows on a log2 scale.
func Level(total int64) int {
	if total <= 0 {
		return 0
	}
	return bits.Len64(uint64(total)) - 1
}

// LeaderboardEntry holds one user and its stats in the leaderboard.
type LeaderboardEntry struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	Avatar   string            `json:"avatar"`
	Stats    GamificationStats `json:"stats"`
}

// stringVal returns the value of a string pointer, or empty string if nil.
func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

// Leaderboard returns the stats of all users.
func (s *GamificationService) Leaderboard() ([]LeaderboardEntry, error) {
	var users []models.User
	if err := s.db.Find(&users).Error; err != nil {
		return nil, err
	}

	entries := make([]LeaderboardEntry, 0, len(users))
	for _, u := range users {
		stats, err := s.Compute(u.ID)
		if err != nil {
			continue
		}
		entries = append(entries, LeaderboardEntry{
			ID:       u.ID,
			Username: u.Username,
			Avatar:   stringVal(u.Avatar),
			Stats:    stats,
		})
	}
	return entries, nil
}
