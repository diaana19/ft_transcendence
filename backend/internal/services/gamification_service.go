// Gamification service — composes a single score from a user's posts, likes
// received on those posts, followers, and following count, then maps that
// score onto a power-of-two level.
package services

import (
	"math/bits"

	"gorm.io/gorm"

	"ft_transcendence/backend/internal/models"
)

// Metric is the per-component shape: a raw count plus its own log2 level so
// each metric carries the same {level, count} envelope as the top-level
// aggregate.
type Metric struct {
	Level int   `json:"level"`
	Count int64 `json:"count"`
}

// asMetric wraps count and pre-computes its floor(log2) level.
func asMetric(count int64) Metric {
	return Metric{Level: Level(count), Count: count}
}

// GamificationStats is the response shape returned by the controller. Each
// component metric is reported as {level, count}; total is the sum of the
// component counts, and the top-level level is floor(log2(total)).
type GamificationStats struct {
	Level     int    `json:"level"`
	Total     int64  `json:"total"`
	Posts     Metric `json:"posts"`
	Likes     Metric `json:"likes"`
	Followers Metric `json:"followers"`
	Following Metric `json:"following"`
}

type GamificationService struct {
	db     *gorm.DB
	friend *FriendService
}

func NewGamificationService(db *gorm.DB, friend *FriendService) *GamificationService {
	return &GamificationService{db: db, friend: friend}
}

// Compute returns total = posts + likes received + followers + following,
// with level derived from total via Level(total). Each underlying query is
// short and indexed; this is intentionally a snapshot, not cached.
func (s *GamificationService) Compute(userID string) (GamificationStats, error) {
	var posts int64
	if err := s.db.Model(&models.Post{}).
		Where("author_id = ?", userID).
		Count(&posts).Error; err != nil {
		return GamificationStats{}, err
	}

	// Likes received = likes on the user's posts plus likes on the user's
	// replies. Both sum the denormalized LikesCount column (maintained by
	// PostService.ReactToPost / ReactToComment) rather than COUNTing the
	// reactions tables — cheaper, and consistent with how the counters are
	// kept. Dislikes are tracked separately and intentionally don't subtract.
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

	followers, err := s.friend.CountFollowers(userID)
	if err != nil {
		return GamificationStats{}, err
	}
	following, err := s.friend.CountFollowing(userID)
	if err != nil {
		return GamificationStats{}, err
	}

	total := posts + likesReceived + followers + following
	return GamificationStats{
		Level:     Level(total),
		Total:     total,
		Posts:     asMetric(posts),
		Likes:     asMetric(likesReceived),
		Followers: asMetric(followers),
		Following: asMetric(following),
	}, nil
}

// Level returns the highest x such that 2^x <= total. Equivalent to
// floor(log2(total)). total <= 0 returns 0 so brand-new accounts start at
// level 0 instead of negative or undefined.
func Level(total int64) int {
	if total <= 0 {
		return 0
	}
	return bits.Len64(uint64(total)) - 1
}

type LeaderboardEntry struct {
	ID       string            `json:"id"`
	Username string            `json:"username"`
	Avatar   string            `json:"avatar"`
	Stats    GamificationStats `json:"stats"`
}

func stringVal(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

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
