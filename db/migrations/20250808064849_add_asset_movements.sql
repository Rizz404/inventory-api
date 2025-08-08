-- +goose Up
CREATE TABLE asset_movements (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NOT NULL,
  from_location_id VARCHAR(26) NULL,
  to_location_id VARCHAR(26) NULL,
  from_user_id VARCHAR(26) NULL,
  to_user_id VARCHAR(26) NULL,
  movement_date TIMESTAMP WITH TIME ZONE NOT NULL,
  moved_by VARCHAR(26) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  FOREIGN KEY (from_location_id) REFERENCES locations(id) ON DELETE
  SET NULL,
    FOREIGN KEY (to_location_id) REFERENCES locations(id) ON DELETE
  SET NULL,
    FOREIGN KEY (from_user_id) REFERENCES users(id) ON DELETE
  SET NULL,
    FOREIGN KEY (to_user_id) REFERENCES users(id) ON DELETE
  SET NULL,
    FOREIGN KEY (moved_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_movements_asset_id ON asset_movements(asset_id);

CREATE INDEX idx_movements_movement_date ON asset_movements(movement_date);

CREATE INDEX idx_movements_to_location_user ON asset_movements(to_location_id, to_user_id);

CREATE TABLE asset_movements_translation (
  id VARCHAR(26) PRIMARY KEY,
  movement_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  notes TEXT NULL,
  UNIQUE (movement_id, lang_code),
  FOREIGN KEY (movement_id) REFERENCES asset_movements(id) ON DELETE CASCADE
);

CREATE INDEX idx_movements_translation_movement_lang ON asset_movements_translation(movement_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_movements_translation_movement_lang;

DROP TABLE IF EXISTS asset_movements_translation;

DROP INDEX IF EXISTS idx_movements_to_location_user;

DROP INDEX IF EXISTS idx_movements_movement_date;

DROP INDEX IF EXISTS idx_movements_asset_id;

DROP TABLE IF EXISTS asset_movements;
