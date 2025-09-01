-- +goose Up
-- +goose StatementBegin
CREATE TABLE subcategories (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) NOT NULL,
    status SMALLINT NOT NULL DEFAULT 1,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES categories (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Index on name for faster search
CREATE INDEX idx_subcategories_name ON subcategories (name);
-- +goose StatementEnd

-- +goose StatementBegin
-- Unique index on slug to avoid duplicates
CREATE UNIQUE INDEX idx_subcategories_slug ON subcategories (slug);
-- +goose StatementEnd

-- +goose StatementBegin
-- Combined index on category_id + name (optional, useful for filtering by category)
CREATE INDEX idx_subcategories_category_name ON subcategories (category_id, name);
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
DROP TABLE IF EXISTS subcategories;
-- +goose StatementEnd