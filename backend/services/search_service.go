package services

import (
	"context"
	"fmt"

	"github.com/Transcendence/models"
)

type userSearcher interface {
	SearchByUsername(ctx context.Context, q string, limit, offset int, sort, order string) ([]models.User, int64, error)
}

type msgSearcher interface {
	SearchByContent(ctx context.Context, userID, q string, limit, offset int, sort, order string) ([]models.Message, int64, error)
}

type SearchService struct {
	userRepo userSearcher
	msgRepo  msgSearcher
}

func NewSearchService(userRepo userSearcher, msgRepo msgSearcher) *SearchService {
	return &SearchService{
		userRepo: userRepo,
		msgRepo:  msgRepo,
	}
}

type SearchResult struct {
	Users    []models.User    `json:"users,omitempty"`
	Messages []models.Message `json:"messages,omitempty"`
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
		messages, total, err := s.msgRepo.SearchByContent(ctx, userID, q, limit, offset, q, order)
		if err != nil {
			return nil, err
		}
		result.Messages = messages
		result.Total = total

	default:
		return nil, fmt.Errorf("Invalid searchtype:", searchType)
	}

	return result, nil
}
