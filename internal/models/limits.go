package models

import "errors"

type LimitAlgorithm string

const (
	LimitAlgorithmFixedWindow   LimitAlgorithm = "fixed_window"
	LimitAlgorithmSlidingWindow LimitAlgorithm = "sliding_window"
)

type RateLimit struct {
	ID         string         `json:"id"`
	Limit      int            `json:"limit"`
	Window     string         `json:"window"`
	Algorithm  LimitAlgorithm `json:"algorithm"`
	KeyPattern string         `json:"key_pattern"`
}

type CreateRateLimitRequest struct {
	RateLimit
}

type RateLimitFilter struct {
	ID         *string         `json:"id"`
	Limit      *int            `json:"limit"`
	Window     *string         `json:"window"`
	Algorithm  *LimitAlgorithm `json:"algorithm"`
	KeyPattern *string         `json:"key_pattern"`
}

func (r *RateLimit) Validate() error {
	if r.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}

	if r.Window == "" {
		return errors.New("window is required")
	}

	if r.Algorithm == "" {
		return errors.New("algorithm is required")
	}

	if !IsValidLimitAlgorithm(string(r.Algorithm)) {
		return errors.New("invalid algorithm")
	}

	if r.KeyPattern == "" {
		return errors.New("key pattern is required")
	}

	return nil
}

func (r *RateLimitFilter) Validate() error {
	if r.ID != nil && *r.ID == "" {
		return errors.New("id is required")
	}

	if r.Limit != nil && *r.Limit <= 0 {
		return errors.New("limit must be greater than 0")
	}

	if r.Window != nil && *r.Window == "" {
		return errors.New("window is required")
	}

	if r.Algorithm != nil && !IsValidLimitAlgorithm(string(*r.Algorithm)) {
		return errors.New("invalid algorithm")
	}

	if r.KeyPattern != nil && *r.KeyPattern == "" {
		return errors.New("key pattern is required")
	}

	return nil
}

func (r *RateLimit) MatchesFilter(filter *RateLimitFilter) bool {
	if filter == nil {
		return true
	}

	if filter.ID != nil && *filter.ID != r.ID {
		return false
	}

	if filter.Limit != nil && *filter.Limit != r.Limit {
		return false
	}

	if filter.Window != nil && *filter.Window != r.Window {
		return false
	}

	if filter.Algorithm != nil && *filter.Algorithm != r.Algorithm {
		return false
	}

	if filter.KeyPattern != nil && *filter.KeyPattern != r.KeyPattern {
		return false
	}

	return true
}

func (r *CreateRateLimitRequest) Validate() error {
	return r.RateLimit.Validate()
}

func IsValidLimitAlgorithm(algorithm string) bool {
	switch algorithm {
	case string(LimitAlgorithmFixedWindow), string(LimitAlgorithmSlidingWindow):
		return true
	}
	return false
}
