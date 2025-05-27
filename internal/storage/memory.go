package storage

import (
	"context"
	"errors"
	"sync"

	"github.com/dnakolan/rate-limiter-service/internal/models"
)

type RateLimitStorage interface {
	Save(ctx context.Context, RateLimit *models.RateLimit) error
	FindAll(ctx context.Context, filter *models.RateLimitFilter) ([]*models.RateLimit, error)
	FindById(ctx context.Context, uid string) (*models.RateLimit, error)
	Delete(ctx context.Context, uid string) error
	Clear(ctx context.Context) error
}

type rateLimitStorage struct {
	sync.RWMutex
	data map[string]*models.RateLimit
}

func NewRateLimitStorage() *rateLimitStorage {
	return &rateLimitStorage{
		data: make(map[string]*models.RateLimit),
	}
}

func (s *rateLimitStorage) Save(ctx context.Context, RateLimit *models.RateLimit) error {
	s.Lock()
	defer s.Unlock()
	s.data[RateLimit.ID] = RateLimit
	return nil
}

func (s *rateLimitStorage) FindAll(ctx context.Context, filter *models.RateLimitFilter) ([]*models.RateLimit, error) {
	s.RLock()
	defer s.RUnlock()
	RateLimits := make([]*models.RateLimit, 0, len(s.data))
	for _, RateLimit := range s.data {
		if filter == nil || RateLimit.MatchesFilter(filter) {
			RateLimits = append(RateLimits, RateLimit)
		}
	}
	return RateLimits, nil
}

func (s *rateLimitStorage) FindById(ctx context.Context, uid string) (*models.RateLimit, error) {
	s.RLock()
	defer s.RUnlock()
	RateLimit, ok := s.data[uid]
	if !ok {
		return nil, errors.New("rateLimit not found")
	}
	return RateLimit, nil
}

func (s *rateLimitStorage) Delete(ctx context.Context, uid string) error {
	s.Lock()
	defer s.Unlock()
	delete(s.data, uid)
	return nil
}

func (s *rateLimitStorage) Clear(ctx context.Context) error {
	s.Lock()
	defer s.Unlock()
	s.data = make(map[string]*models.RateLimit)
	return nil
}
