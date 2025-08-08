-- +goose Up
CREATE TABLE locations (
  id VARCHAR(26) PRIMARY KEY,
  location_code VARCHAR(20) UNIQUE NOT NULL,
  building VARCHAR(100) NULL,
  floor VARCHAR(20) NULL
);

CREATE INDEX idx_locations_building_floor ON locations(building, floor);

CREATE TABLE locations_translation (
  id VARCHAR(26) PRIMARY KEY,
  location_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  location_name VARCHAR(100) NOT NULL,
  UNIQUE (location_id, lang_code),
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

CREATE INDEX idx_locations_translation_location_lang ON locations_translation(location_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_locations_translation_location_lang;

DROP TABLE IF EXISTS locations_translation;

DROP INDEX IF EXISTS idx_locations_building_floor;

DROP TABLE IF EXISTS locations;
