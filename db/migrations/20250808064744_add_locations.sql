-- +goose Up
CREATE TABLE locations (
  id VARCHAR(26) PRIMARY KEY,
  location_code VARCHAR(20) UNIQUE NOT NULL,
  building VARCHAR(100) NULL,
  floor VARCHAR(20) NULL,
  latitude DECIMAL(11,8) NULL,
  longitude DECIMAL(11,8) NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_locations_building_floor ON locations(building, floor);

CREATE INDEX idx_locations_coordinates ON locations(latitude, longitude);

CREATE TABLE location_translations (
  id VARCHAR(26) PRIMARY KEY,
  location_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  location_name VARCHAR(100) NOT NULL,
  UNIQUE (location_id, lang_code),
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE CASCADE
);

CREATE INDEX idx_location_translations_location_lang ON location_translations(location_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_location_translations_location_lang;

DROP TABLE IF EXISTS location_translations;

DROP INDEX IF EXISTS idx_locations_coordinates;

DROP INDEX IF EXISTS idx_locations_building_floor;

DROP TABLE IF EXISTS locations;
