-- +goose Up
-- +goose StatementBegin
CREATE TABLE models (
                        id BIGSERIAL PRIMARY KEY,
                        brand_id BIGINT,
                        name VARCHAR(255) NOT NULL,
                        release_date DATE,
                        is_upcoming BOOLEAN DEFAULT FALSE,
                        is_limited BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для поиска моделей по имени
CREATE INDEX idx_models_name ON models (name);

-- Индекс для фильтрации по brand_id (связь с брендом)
CREATE INDEX idx_models_brand_id ON models (brand_id);

-- Индекс для фильтрации по дате релиза
CREATE INDEX idx_models_release_date ON models (release_date);

-- Индекс для фильтрации по флагу is_upcoming
CREATE INDEX idx_models_is_upcoming ON models (is_upcoming);

-- Индекс для фильтрации по флагу is_limited
CREATE INDEX idx_models_is_limited ON models (is_limited);

-- Индекс для временных запросов по дате создания
CREATE INDEX idx_models_created_at ON models (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_models_created_at;
DROP INDEX IF EXISTS idx_models_is_limited;
DROP INDEX IF EXISTS idx_models_is_upcoming;
DROP INDEX IF EXISTS idx_models_release_date;
DROP INDEX IF EXISTS idx_models_brand_id;
DROP INDEX IF EXISTS idx_models_name;
DROP TABLE IF EXISTS models;
-- +goose StatementEnd
