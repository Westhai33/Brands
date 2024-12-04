package model

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// SoftDelete мягко удаляет модель
func (s *Service) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.SoftDelete")
	defer span.Finish()

	err := s.repo.SoftDelete(ctx, id)
	if err != nil {
		span.SetTag("error", true)
		return err
	}
	return nil
}

// Restore восстанавливает мягко удалённую модель
func (s *Service) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Restore")
	defer span.Finish()

	err := s.repo.Restore(ctx, id)
	if err != nil {
		span.SetTag("error", true)
		return err
	}
	return nil
}
