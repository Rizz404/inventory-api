-- +goose Up
-- +goose StatementBegin
CREATE TABLE asset_images (
  id TEXT PRIMARY KEY,
  asset_id TEXT NOT NULL,
  image_id TEXT NOT NULL,
  display_order INTEGER DEFAULT 0,
  is_primary BOOLEAN DEFAULT false,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_asset_images_asset FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  CONSTRAINT fk_asset_images_image FOREIGN KEY (image_id) REFERENCES images(id) ON DELETE CASCADE,
  CONSTRAINT uk_asset_image UNIQUE(asset_id, image_id)
);

CREATE INDEX idx_asset_images_asset_id ON asset_images(asset_id);

CREATE INDEX idx_asset_images_image_id ON asset_images(image_id);

CREATE INDEX idx_asset_images_is_primary ON asset_images(is_primary);

CREATE INDEX idx_asset_images_display_order ON asset_images(display_order);

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_asset_images_display_order;

DROP INDEX IF EXISTS idx_asset_images_is_primary;

DROP INDEX IF EXISTS idx_asset_images_image_id;

DROP INDEX IF EXISTS idx_asset_images_asset_id;

DROP TABLE IF EXISTS asset_images;

-- +goose StatementEnd
