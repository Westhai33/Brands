package dto

import "time"

type Model struct {
	ID          int64      `json:"id"`           // BIGSERIAL
	BrandID     int64      `json:"brand_id"`     // Ссылка на бренд
	Name        string     `json:"name"`         // Название модели
	ReleaseDate *time.Time `json:"release_date"` // Дата релиза
	IsUpcoming  bool       `json:"is_upcoming"`  // Флаг "Скоро"
	IsLimited   bool       `json:"is_limited"`   // Флаг ограниченного выпуска
	CreatedAt   time.Time  `json:"created_at"`   // Время создания
	UpdatedAt   time.Time  `json:"updated_at"`   // Время обновления
}
