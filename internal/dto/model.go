package dto

import "time"

type Model struct {
	ID          int64     `json:"id"`
	BrandID     int64     `json:"brand_id"`
	Name        string    `json:"name"`         // Название модели
	Link        string    `json:"link"`         // Модель на английском
	ReleaseDate time.Time `json:"release_date"` // Дата релиза
	IsUpcoming  bool      `json:"is_upcoming"`  // Флаг "Скоро"
	IsLimited   bool      `json:"is_limited"`   // Флаг ограниченного выпуска
	IsDeleted   bool      `json:"is_deleted"`   // Флаг удаления

	CreatedAt time.Time `json:"created_at"` // Время создания
	UpdatedAt time.Time `json:"updated_at"` // Время обновления
}
