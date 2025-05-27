package storage

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/dnakolan/rate-limiter-service/internal/models"
)

func TestRateLimitStorage_Save(t *testing.T) {
	storage := NewRateLimitStorage()
	ctx := context.Background()

	rateLimit := &models.RateLimit{
		ID:         "api-users",
		Limit:      200,
		Window:     "1m",
		Algorithm:  models.LimitAlgorithmSlidingWindow,
		KeyPattern: "user:{{user_id}}",
	}

	err := storage.Save(ctx, rateLimit)
	require.NoError(t, err)

	// Verify the rateLimit was saved
	saved, err := storage.FindById(ctx, rateLimit.ID)
	require.NoError(t, err)
	assert.Equal(t, *rateLimit, *saved)
}

func TestRateLimitStorage_FindById(t *testing.T) {
	storage := NewRateLimitStorage()
	ctx := context.Background()

	rateLimit := &models.RateLimit{
		ID:         "api-users",
		Limit:      200,
		Window:     "1m",
		Algorithm:  models.LimitAlgorithmSlidingWindow,
		KeyPattern: "user:{{user_id}}",
	}

	// Save a rateLimit first
	err := storage.Save(ctx, rateLimit)
	require.NoError(t, err)

	tests := []struct {
		name          string
		id            string
		expectError   bool
		expectedError string
	}{
		{
			name:        "successful retrieval",
			id:          rateLimit.ID,
			expectError: false,
		},
		{
			name:          "rateLimit not found",
			id:            "nonexistent",
			expectError:   true,
			expectedError: "rateLimit not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindById(ctx, tt.id)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, found)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, found)
				assert.Equal(t, *rateLimit, *found)
			}
		})
	}
}

func TestRateLimitStorage_FindAll(t *testing.T) {
	storage := NewRateLimitStorage()
	ctx := context.Background()

	// Create test rateLimits
	rateLimits := []*models.RateLimit{
		&models.RateLimit{
			ID:         "api-users",
			Limit:      200,
			Window:     "1m",
			Algorithm:  models.LimitAlgorithmSlidingWindow,
			KeyPattern: "user:{{user_id}}",
		},
		&models.RateLimit{
			ID:         "api-users-2",
			Limit:      300,
			Window:     "1m",
			Algorithm:  models.LimitAlgorithmFixedWindow,
			KeyPattern: "user:{{user_id}}",
		},
		&models.RateLimit{
			ID:         "api-users-3",
			Limit:      400,
			Window:     "1m",
			Algorithm:  models.LimitAlgorithmSlidingWindow,
			KeyPattern: "user:{{user_id}}",
		},
	}

	// Save all rateLimits
	for _, w := range rateLimits {
		err := storage.Save(ctx, w)
		require.NoError(t, err)
	}

	tests := []struct {
		name           string
		filter         *models.RateLimitFilter
		expectedCount  int
		expectedIds    []string
		expectedFilter func(*models.RateLimit) bool
	}{
		{
			name:          "no filter",
			filter:        nil,
			expectedCount: 3,
			expectedIds:   []string{"api-users", "api-users-2", "api-users-3"},
		},
		{
			name: "filter by user id",
			filter: &models.RateLimitFilter{
				ID: stringPtr("api-users-2"),
			},
			expectedCount: 1,
			expectedIds:   []string{"api-users-2"},
			expectedFilter: func(w *models.RateLimit) bool {
				return w.ID == "api-users-2"
			},
		},
		{
			name: "filter by algorithm",
			filter: &models.RateLimitFilter{
				Algorithm: algorithmPtr("fixed_window"),
			},
			expectedCount: 1,
			expectedIds:   []string{"api-users-2"},
			expectedFilter: func(w *models.RateLimit) bool {
				return w.Algorithm == models.LimitAlgorithmFixedWindow
			},
		},
		{
			name: "filter with no matches",
			filter: &models.RateLimitFilter{
				ID: stringPtr("nonexistent"),
			},
			expectedCount: 0,
			expectedIds:   []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			found, err := storage.FindAll(ctx, tt.filter)
			require.NoError(t, err)
			assert.Len(t, found, tt.expectedCount)

			// Verify the names of returned rateLimits
			names := make([]string, len(found))
			for i, w := range found {
				names[i] = w.ID
			}
			assert.ElementsMatch(t, tt.expectedIds, names)

			// If there's a specific filter function, verify each rateLimit matches it
			if tt.expectedFilter != nil {
				for _, w := range found {
					assert.True(t, tt.expectedFilter(w))
				}
			}
		})
	}
}

// Helper functions to create pointers
func float64Ptr(v float64) *float64 {
	return &v
}

func stringPtr(v string) *string {
	return &v
}

func algorithmPtr(v string) *models.LimitAlgorithm {
	if v == "" {
		return nil
	}
	algorithm := models.LimitAlgorithm(v)
	return &algorithm
}
