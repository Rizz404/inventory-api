-- +goose Up
CREATE TYPE scan_method_type AS ENUM ('DATA_MATRIX', 'MANUAL_INPUT');

CREATE TYPE scan_result_type AS ENUM ('Success', 'Invalid ID', 'Asset Not Found');

CREATE TABLE scan_logs (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NULL,
  scanned_value VARCHAR(255) NOT NULL,
  scan_method scan_method_type NOT NULL,
  scanned_by VARCHAR(26) NOT NULL,
  scan_timestamp TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  scan_location_lat DECIMAL(11, 8) NULL,
  scan_location_lng DECIMAL(11, 8) NULL,
  scan_result scan_result_type NOT NULL,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE
  SET NULL,
    FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_scan_logs_scan_timestamp ON scan_logs(scan_timestamp);

CREATE INDEX idx_scan_logs_scanned_by ON scan_logs(scanned_by);

CREATE INDEX idx_scan_logs_result ON scan_logs(scan_result);

CREATE INDEX idx_scan_logs_location ON scan_logs(scan_location_lat, scan_location_lng);

-- +goose Down
DROP INDEX IF EXISTS idx_scan_logs_location;

DROP INDEX IF EXISTS idx_scan_logs_result;

DROP INDEX IF EXISTS idx_scan_logs_scanned_by;

DROP INDEX IF EXISTS idx_scan_logs_scan_timestamp;

DROP TABLE IF EXISTS scan_logs;

DROP TYPE IF EXISTS scan_result_type;

DROP TYPE IF EXISTS scan_method_type;
