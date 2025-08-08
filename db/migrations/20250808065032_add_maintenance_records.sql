-- +goose Up
CREATE TABLE maintenance_records (
  id VARCHAR(26) PRIMARY KEY,
  schedule_id VARCHAR(26) NULL,
  asset_id VARCHAR(26) NOT NULL,
  maintenance_date DATE NOT NULL,
  performed_by_user VARCHAR(26) NULL,
  performed_by_vendor VARCHAR(150) NULL,
  actual_cost DECIMAL(12, 2) NULL,
  FOREIGN KEY (schedule_id) REFERENCES maintenance_schedules(id) ON DELETE
  SET NULL,
    FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
    FOREIGN KEY (performed_by_user) REFERENCES users(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_maintenance_records_asset_id ON maintenance_records(asset_id);

CREATE INDEX idx_maintenance_records_date ON maintenance_records(maintenance_date);

CREATE TABLE maintenance_records_translation (
  id VARCHAR(26) PRIMARY KEY,
  record_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  title VARCHAR(200) NOT NULL,
  notes TEXT NULL,
  UNIQUE (record_id, lang_code),
  FOREIGN KEY (record_id) REFERENCES maintenance_records(id) ON DELETE CASCADE
);

CREATE INDEX idx_records_translation_record_lang ON maintenance_records_translation(record_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_records_translation_record_lang;

DROP TABLE IF EXISTS maintenance_records_translation;

DROP INDEX IF EXISTS idx_maintenance_records_date;

DROP INDEX IF EXISTS idx_maintenance_records_asset_id;

DROP TABLE IF EXISTS maintenance_records;
