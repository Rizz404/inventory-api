-- +goose Up
CREATE TYPE maintenance_schedule_type AS ENUM ('Preventive', 'Corrective');

CREATE TYPE schedule_status AS ENUM ('Scheduled', 'Completed', 'Cancelled');

CREATE TABLE maintenance_schedules (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NOT NULL,
  maintenance_type maintenance_schedule_type NOT NULL,
  scheduled_date DATE NOT NULL,
  frequency_months INT NULL,
  status schedule_status DEFAULT 'Scheduled',
  created_by VARCHAR(26) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT
);

CREATE INDEX idx_maintenance_schedules_asset_id ON maintenance_schedules(asset_id);

CREATE INDEX idx_maintenance_schedules_status ON maintenance_schedules(status);

CREATE INDEX idx_maintenance_schedules_scheduled_date ON maintenance_schedules(scheduled_date);

CREATE TABLE maintenance_schedule_translations (
  id VARCHAR(26) PRIMARY KEY,
  schedule_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NULL,
  UNIQUE (schedule_id, lang_code),
  FOREIGN KEY (schedule_id) REFERENCES maintenance_schedules(id) ON DELETE CASCADE
);

CREATE INDEX idx_schedule_translations_schedule_lang ON maintenance_schedule_translations(schedule_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_schedule_translations_schedule_lang;

DROP TABLE IF EXISTS maintenance_schedule_translations;

DROP INDEX IF EXISTS idx_maintenance_schedules_scheduled_date;

DROP INDEX IF EXISTS idx_maintenance_schedules_status;

DROP INDEX IF EXISTS idx_maintenance_schedules_asset_id;

DROP TABLE IF EXISTS maintenance_schedules;

DROP TYPE IF EXISTS schedule_status;

DROP TYPE IF EXISTS maintenance_schedule_type;
