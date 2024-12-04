package utils

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/valyala/fasthttp"
)

// ExtractUUIDFromPath извлекает и парсит UUID из пути запроса.
// Возвращает uuid.UUID и nil, если успешно.
// Возвращает ошибку, если extraction или парсинг не удался.
func ExtractUUIDFromPath(ctx *fasthttp.RequestCtx, key string) (uuid.UUID, error) {
	idStr, ok := ctx.UserValue(key).(string)
	if !ok || idStr == "" {
		return uuid.Nil, fmt.Errorf("missing %s in path", key)
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid %s format: %w", key, err)
	}

	return id, nil
}
