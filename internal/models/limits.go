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
