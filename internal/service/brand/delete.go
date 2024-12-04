package brand

import (
	"context"
	"github.com/opentracing/opentracing-go"
)

// SoftDelete мягко удаляет бренд
func (s *BrandService) SoftDelete(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.SoftDelete")
	defer span.Finish()

	err := s.repo.SoftDelete(ctx, id)
	if err != nil {

		return err
	}
	return nil
}

// Restore восстанавливает мягко удалённый бренд
func (s *BrandService) Restore(ctx context.Context, id int64) error {
	span, ctx := opentracing.StartSpanFromContext(ctx, "BrandService.Restore")
	defer span.Finish()

	err := s.repo.Restore(ctx, id)
	if err != nil {

		return err
	}
	return nil
}
