-- +goose Up
CREATE TYPE maintenance_type AS ENUM (
  'Preventive',
  'Corrective',
  'Inspection',
  'Calibration'
);

CREATE TYPE schedule_state AS ENUM ('Active', 'Paused', 'Stopped', 'Completed');

CREATE TYPE interval_unit AS ENUM (
  'Minutes',
  'Hours',
  'Days',
  'Weeks',
  'Months',
  'Years'
);

CREATE TABLE maintenance_schedules (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NOT NULL,
  maintenance_type maintenance_type NOT NULL,
  is_recurring BOOLEAN DEFAULT FALSE,
  -- Flexible interval: "setiap X unit"
  -- Contoh: interval_value=15, interval_unit='Days' = setiap 15 hari
  -- Contoh: interval_value=2, interval_unit='Months' = setiap 2 bulan
  interval_value INT NULL,
  -- berapa kali interval
  interval_unit interval_unit NULL,
  -- satuan waktu
  -- Waktu spesifik untuk eksekusi (format: HH:MM:SS)
  -- Null = kapan aja dalam hari tersebut
  scheduled_time TIME NULL,
  -- Tanggal maintenance berikutnya (auto-update oleh cron)
  next_scheduled_date TIMESTAMP WITH TIME ZONE NOT NULL,
  -- Tanggal maintenance terakhir kali dilakukan
  last_executed_date TIMESTAMP WITH TIME ZONE NULL,
  -- State management
  state schedule_state DEFAULT 'Active',
  -- Auto-completion: jika TRUE, status jadi 'Completed' setelah 1x maintenance
  auto_complete BOOLEAN DEFAULT FALSE,
  -- Estimasi biaya (untuk budgeting)
  estimated_cost DECIMAL(12, 2) NULL,
  created_by VARCHAR(26) NOT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  FOREIGN KEY (created_by) REFERENCES users(id) ON DELETE RESTRICT,
  -- Validasi: jika recurring, harus ada interval
  CONSTRAINT check_recurring_interval CHECK (
    (is_recurring = FALSE)
    OR (
      is_recurring = TRUE
      AND interval_value IS NOT NULL
      AND interval_unit IS NOT NULL
    )
  )
);

CREATE INDEX idx_maintenance_schedules_asset_id ON maintenance_schedules(asset_id);

CREATE INDEX idx_maintenance_schedules_state ON maintenance_schedules(state);

CREATE INDEX idx_maintenance_schedules_next_date ON maintenance_schedules(next_scheduled_date);

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

DROP INDEX IF EXISTS idx_maintenance_schedules_next_date;

DROP INDEX IF EXISTS idx_maintenance_schedules_state;

DROP INDEX IF EXISTS idx_maintenance_schedules_asset_id;

DROP TABLE IF EXISTS maintenance_schedules;

DROP TYPE IF EXISTS interval_unit;

DROP TYPE IF EXISTS schedule_state;

DROP TYPE IF EXISTS maintenance_type;
