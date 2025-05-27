package services

import (
	"context"

	"github.com/dnakolan/rate-limiter-service/internal/models"
	"github.com/dnakolan/rate-limiter-service/internal/storage"
)

type LimitsService interface {
	CreateRateLimit(ctx context.Context, req *models.RateLimit) error
}

type limitsService struct {
	storage storage.RateLimitStorage
}

func NewLimitsService(storage storage.RateLimitStorage) *limitsService {
	return &limitsService{
		storage: storage,
	}
}

func (s *limitsService) CreateRateLimit(ctx context.Context, req *models.RateLimit) error {
	return s.storage.Save(ctx, req)
}
