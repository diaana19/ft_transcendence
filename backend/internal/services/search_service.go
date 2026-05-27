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
	Users    []models.User    `json:"users,omitempty"`
	Messages []models.Message `json:"messages,omitempty"`
	Posts    []models.Post    `json:"posts,omitempty"`
	Total    int64            `json:"total"`
	Page     int              `json:"page"`
	Limit    int              `json:"limit"`
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
		result.Users = users
		result.Total = total
	case "message":
		messages, total, err := s.msgRepo.SearchByContent(ctx, userID, q, limit, offset, sort, order)
		if err != nil {
			return nil, err
		}
		result.Messages = messages
		result.Total = total
	case "post":
		posts, total, err := s.postRepo.SearchByPost(ctx, q, limit, offset, sort, order)
		if err != nil {
			return nil, err
		}
		result.Posts = posts
		result.Total = total
	case "all", "":
		users, total, err := s.userRepo.SearchByUsername(ctx, q, limit, 0, "username", "asc")
		if err != nil {
			return nil, err
		}
		result.Users = users
		result.Total = total
		posts, total, err := s.postRepo.SearchByPost(ctx, q, limit, 0, "created_at", "desc")
		if err != nil {
			return nil, err
		}
		result.Posts = posts
		result.Total += total

	default:
		return nil, fmt.Errorf("invalid search type: %s", searchType)
	}

	return result, nil
}
