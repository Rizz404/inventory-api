-- +goose Up
-- +goose StatementBegin
CREATE TABLE images (
  id TEXT PRIMARY KEY,
  image_url TEXT NOT NULL,
  public_id TEXT,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_images_public_id ON images(public_id);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_images_public_id;

DROP TABLE IF EXISTS images;

-- +goose StatementEnd
