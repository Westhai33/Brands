package dto

import (
	"github.com/google/uuid"
	"time"
)

// Brand представляет сущность бренда
type Brand struct {
	ID            uuid.UUID `json:"id"`
	Name          string    `json:"name"`            // Название бренда
	Link          string    `json:"link"`            // Бренд на английском
	Description   string    `json:"description"`     // Описание/история бренда
	LogoURL       string    `json:"logo_url"`        // URL логотипа
	CoverImageURL string    `json:"cover_image_url"` // URL обложки
	FoundedYear   int       `json:"founded_year"`    // Год основания
	OriginCountry string    `json:"origin_country"`  // Страна происхождения
	Popularity    int       `json:"popularity"`      // Индекс популярности
	IsPremium     bool      `json:"is_premium"`      // Флаг премиального бренда
	IsUpcoming    bool      `json:"is_upcoming"`     // Флаг "Скоро"
	IsDeleted     bool      `json:"is_deleted"`      // Флаг удаления

	CreatedAt time.Time `json:"created_at"` // Время создания
	UpdatedAt time.Time `json:"updated_at"` // Время обновления
}
