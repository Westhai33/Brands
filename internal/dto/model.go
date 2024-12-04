package dto

import (
	"github.com/google/uuid"
	"time"
)

type Model struct {
	ID          uuid.UUID `json:"id"`
	BrandID     uuid.UUID `json:"brand_id"`
	Name        string    `json:"name"`         // Название модели
	ReleaseDate time.Time `json:"release_date"` // Дата релиза
	IsUpcoming  bool      `json:"is_upcoming"`  // Флаг "Скоро"
	IsLimited   bool      `json:"is_limited"`   // Флаг ограниченного выпуска
	IsDeleted   bool      `json:"is_deleted"`   // Флаг удаления

	CreatedAt time.Time `json:"created_at"` // Время создания
	UpdatedAt time.Time `json:"updated_at"` // Время обновления
}
