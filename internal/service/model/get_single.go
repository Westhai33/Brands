package model

import (
	"Brands/internal/dto"
	"context"
	"github.com/opentracing/opentracing-go"
)

func (s *Service) GetByID(ctx context.Context, id int64) (*dto.Model, error) {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.GetByID")
	defer span.Finish()

	model, err := s.repo.GetByID(ctx, id)
	if err != nil {
		span.SetTag("error", true)
		return nil, err
	}
	return model, nil
}
