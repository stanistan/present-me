package api

import (
	"context"
)

type ErrSource struct {
	Err error
}

var _ Source = &ErrSource{}

func (s *ErrSource) GetReview(_ context.Context) (Review, error) {
	return Review{}, s.Err
}
