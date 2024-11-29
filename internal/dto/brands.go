package dto

import (
	"encoding/json"
	"fmt"
	"time"
)

// Brand представляет сущность бренда
type Brand struct {
	ID            int64     `json:"id"`              // BIGSERIAL
	Name          string    `json:"name"`            // NOT NULL
	Link          string    `json:"link"`            // Может быть пустым
	Description   string    `json:"description"`     // Описание/история бренда
	LogoURL       string    `json:"logo_url"`        // URL логотипа
	CoverImageURL string    `json:"cover_image_url"` // URL обложки
	FoundedYear   int       `json:"founded_year"`    // Год основания
	OriginCountry string    `json:"origin_country"`  // Страна происхождения
	Popularity    int       `json:"popularity"`      // Индекс популярности
	IsPremium     bool      `json:"is_premium"`      // Флаг премиального бренда
	IsUpcoming    bool      `json:"is_upcoming"`     // Флаг "Скоро"
	IsDeleted     bool      `json:"is_deleted"`      // Флаг удаления
	CreatedAt     time.Time `json:"created_at"`      // Время создания
	UpdatedAt     time.Time `json:"updated_at"`      // Время обновления
}

func (b *Brand) ToJSON() (string, error) {
	data, err := json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("ошибка преобразования в JSON: %w", err)
	}
	return string(data), nil
}
