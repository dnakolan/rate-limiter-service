package services

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dnakolan/rate-limiter-service/internal/models"
)

type RateLimiterEngine interface {
	CheckRateLimit(ctx context.Context, userID string, rule *models.RateLimit) bool
	ApplyRateLimit(ctx context.Context, userID string, rule *models.RateLimit) error
}

type rateLimiterEngine struct {
	mu      sync.RWMutex
	windows map[string][]time.Time
}

func NewRateLimiterEngine() *rateLimiterEngine {
	return &rateLimiterEngine{
		windows: make(map[string][]time.Time),
	}
}

func (r *rateLimiterEngine) ApplyRateLimit(ctx context.Context, userID string, rule *models.RateLimit) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	key := fmt.Sprintf("%s:%s", rule.ID, userID)
	r.windows[key] = append(r.windows[key], time.Now())

	return nil
}

func (r *rateLimiterEngine) CheckRateLimit(ctx context.Context, userID string, rule *models.RateLimit) bool {
	window, err := time.ParseDuration(rule.Window)
	if err != nil {
		slog.Error("failed to parse window duration", "error", err)
		return true
	}

	key := fmt.Sprintf("%s:%s", rule.ID, userID)
	return r.slidingWindowCheck(key, rule.Limit, window)
}

func (r *rateLimiterEngine) slidingWindowCheck(key string, limit int, window time.Duration) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()
	cutoff := now.Add(-window)

	// Remove old requests outside the window
	timestamps := r.windows[key]
	valid := []time.Time{}
	for _, t := range timestamps {
		if t.After(cutoff) { // keep requests within window
			valid = append(valid, t)
		}
	}

	// Check if we can add another request
	if len(valid) < limit {
		valid = append(valid, now)
		r.windows[key] = valid
		return true
	}

	r.windows[key] = valid // update with cleaned list
	return false
}
