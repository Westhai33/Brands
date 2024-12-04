package model

import (
	"context"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// SoftDelete мягко удаляет модель
func (s *ModelService) SoftDelete(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.SoftDelete")
	defer span.Finish()

	err := s.repo.SoftDelete(ctx, id)
	if err != nil {

		return err
	}
	return nil
}

// Restore восстанавливает мягко удалённую модель
func (s *ModelService) Restore(ctx context.Context, id uuid.UUID) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "ModelService.Restore")
	defer span.Finish()

	err := s.repo.Restore(ctx, id)
	if err != nil {

		return err
	}
	return nil
}
