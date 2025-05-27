package services

import (
	"context"

	"github.com/dnakolan/rate-limiter-service/internal/models"
	"github.com/dnakolan/rate-limiter-service/internal/storage"
)

type LimitsService interface {
	CreateRateLimit(ctx context.Context, req *models.RateLimit) error
	GetRateLimit(ctx context.Context, id string) (*models.RateLimit, error)
	DeleteRateLimit(ctx context.Context, id string) error
	UpdateRateLimit(ctx context.Context, req *models.RateLimit) error

	CheckRateLimit(ctx context.Context, ruleId, userId string) (bool, error)
	ApplyRateLimit(ctx context.Context, ruleId, userId string) error
}

type limitsService struct {
	storage storage.RateLimitStorage
	engine  RateLimiterEngine
}

func NewLimitsService(rateLimitStorage storage.RateLimitStorage) *limitsService {
	return &limitsService{
		storage: rateLimitStorage,
		engine:  NewRateLimiterEngine(),
	}
}

func (s *limitsService) CreateRateLimit(ctx context.Context, req *models.RateLimit) error {
	return s.storage.Save(ctx, req)
}

func (s *limitsService) GetRateLimit(ctx context.Context, id string) (*models.RateLimit, error) {
	return s.storage.FindById(ctx, id)
}

func (s *limitsService) DeleteRateLimit(ctx context.Context, id string) error {
	return s.storage.Delete(ctx, id)
}

func (s *limitsService) UpdateRateLimit(ctx context.Context, req *models.RateLimit) error {
	return s.storage.Save(ctx, req)
}

func (s *limitsService) CheckRateLimit(ctx context.Context, userId, ruleId string) (bool, error) {
	rule, err := s.storage.FindById(ctx, ruleId)
	if err != nil {
		return true, err
	}

	return s.engine.CheckRateLimit(ctx, userId, rule), nil
}

func (s *limitsService) ApplyRateLimit(ctx context.Context, ruleId, userId string) error {
	rule, err := s.storage.FindById(ctx, ruleId)
	if err != nil {
		return err
	}

	return s.engine.ApplyRateLimit(ctx, userId, rule)
}
