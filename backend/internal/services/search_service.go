package services

import (
	"context"
	"fmt"

	"ft_transcendence/backend/internal/models"
)

type userSearcher interface {
	SearchByUsername(ctx context.Context, q string, limit, offset int, sort, order string) ([]models.User, int64, error)
}

type msgSearcher interface {
	SearchByContent(
		ctx context.Context, userID, q string, limit, offset int, sort, order string,
	) ([]models.Message, int64, error)
}

type postSearcher interface {
	SearchByPost(ctx context.Context, q string, limit, offset int, sort, order string) ([]models.Post, int64, error)
}

type SearchService struct {
	userRepo userSearcher
	msgRepo  msgSearcher
	postRepo postSearcher
}

func NewSearchService(userRepo userSearcher, msgRepo msgSearcher, postRepo postSearcher) *SearchService {
	return &SearchService{
		userRepo: userRepo,
		msgRepo:  msgRepo,
		postRepo: postRepo,
	}
}

type SearchResult struct {
	Users    []models.UserResponse    `json:"users,omitempty"`
	Messages []models.MessageResponse `json:"messages,omitempty"`
	Posts    []models.PostResponse    `json:"posts,omitempty"`
	Total    int64                    `json:"total"`
	Page     int                      `json:"page"`
	Limit    int                      `json:"limit"`
}

// toUserResponses, toMessageResponses, toPostResponses convert raw model
// slices coming out of the repositories into the DTO shapes the rest of the
// API uses. Kept private to the service — repositories deal in models,
// controllers see only DTOs.
func toUserResponses(users []models.User) []models.UserResponse {
	out := make([]models.UserResponse, len(users))
	for i := range users {
		out[i] = users[i].ToResponse()
	}
	return out
}

func toMessageResponses(msgs []models.Message) []models.MessageResponse {
	out := make([]models.MessageResponse, len(msgs))
	for i := range msgs {
		out[i] = msgs[i].ToResponse()
	}
	return out
}

func toPostResponses(posts []models.Post) []models.PostResponse {
	out := make([]models.PostResponse, len(posts))
	for i := range posts {
		out[i] = posts[i].ToResponse()
	}
	return out
}

func (s *SearchService) Search(
	ctx context.Context,
	userID string,
	q, searchType, sort, order string,
	page, limit int,
) (*SearchResult, error) {
	offset := (page - 1) * limit
	result := &SearchResult{Page: page, Limit: limit}

	switch searchType {
	case "user":
		users, total, err := s.userRepo.SearchByUsername(ctx, q, limit, offset, sort, order)
		if err != nil {
			return nil, err
		}
		result.Users = toUserResponses(users)
		result.Total = total
	case "message":
		messages, total, err := s.msgRepo.SearchByContent(ctx, userID, q, limit, offset, sort, order)
		if err != nil {
			return nil, err
		}
		result.Messages = toMessageResponses(messages)
		result.Total = total
	case "post":
		posts, total, err := s.postRepo.SearchByPost(ctx, q, limit, offset, sort, order)
		if err != nil {
			return nil, err
		}
		result.Posts = toPostResponses(posts)
		result.Total = total
	case "all", "":
		// "all" runs the three per-type queries with the caller-supplied
		// pagination applied to each section independently. Each section keeps
		// its natural default ordering (users by username asc, posts and
		// messages by created_at desc) — the global "sort" parameter is not
		// meaningful across heterogeneous result kinds. Total is the sum of
		// matching rows across all three types, so a client can show
		// "N total matches" without paginating combined results.
		users, uTotal, err := s.userRepo.SearchByUsername(ctx, q, limit, offset, "username", "asc")
		if err != nil {
			return nil, err
		}
		posts, pTotal, err := s.postRepo.SearchByPost(ctx, q, limit, offset, "created_at", "desc")
		if err != nil {
			return nil, err
		}
		messages, mTotal, err := s.msgRepo.SearchByContent(ctx, userID, q, limit, offset, "created_at", "desc")
		if err != nil {
			return nil, err
		}
		result.Users = toUserResponses(users)
		result.Posts = toPostResponses(posts)
		result.Messages = toMessageResponses(messages)
		result.Total = uTotal + pTotal + mTotal

	default:
		return nil, fmt.Errorf("invalid search type: %s", searchType)
	}

	return result, nil
}
