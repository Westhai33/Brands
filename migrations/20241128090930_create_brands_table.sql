-- +goose Up
-- +goose StatementBegin
CREATE TABLE brands (
                        id uuid NOT NULL PRIMARY KEY,
                        name VARCHAR(255) NOT NULL,
                        link VARCHAR(255),
                        description TEXT,
                        logo_url VARCHAR(255),
                        cover_image_url VARCHAR(255),
                        founded_year INTEGER,
                        origin_country VARCHAR(100),
                        popularity INTEGER DEFAULT 0,
                        is_premium BOOLEAN DEFAULT FALSE,
                        is_upcoming BOOLEAN DEFAULT FALSE,
                        is_deleted BOOLEAN DEFAULT FALSE,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс для поиска брендов по имени
CREATE INDEX idx_brands_name ON brands (name);

-- Индекс для фильтрации по стране происхождения
CREATE INDEX idx_brands_origin_country ON brands (origin_country);

-- Индекс для сортировки и фильтрации по популярности
CREATE INDEX idx_brands_popularity ON brands (popularity);

-- Составной индекс для фильтрации по премиальным брендам и популярности
CREATE INDEX idx_brands_premium_popularity ON brands (is_premium, popularity);

-- Индекс для временных запросов по дате создания
CREATE INDEX idx_brands_created_at ON brands (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_brands_created_at;
DROP INDEX IF EXISTS idx_brands_premium_popularity;
DROP INDEX IF EXISTS idx_brands_popularity;
DROP INDEX IF EXISTS idx_brands_origin_country;
DROP INDEX IF EXISTS idx_brands_name;
DROP TABLE IF EXISTS brands;
-- +goose StatementEnd
