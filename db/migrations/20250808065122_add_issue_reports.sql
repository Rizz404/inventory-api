-- +goose Up
CREATE TYPE issue_priority AS ENUM ('Low', 'Medium', 'High', 'Critical');

CREATE TYPE issue_status AS ENUM ('Open', 'In Progress', 'Resolved', 'Closed');

CREATE TABLE issue_reports (
  id VARCHAR(26) PRIMARY KEY,
  asset_id VARCHAR(26) NOT NULL,
  reported_by VARCHAR(26) NOT NULL,
  reported_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  issue_type VARCHAR(50) NOT NULL,
  priority issue_priority DEFAULT 'Medium',
  status issue_status DEFAULT 'Open',
  resolved_date TIMESTAMP WITH TIME ZONE NULL,
  resolved_by VARCHAR(26) NULL,
  FOREIGN KEY (asset_id) REFERENCES assets(id) ON DELETE CASCADE,
  FOREIGN KEY (reported_by) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (resolved_by) REFERENCES users(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_issue_reports_asset_id ON issue_reports(asset_id);

CREATE INDEX idx_issue_reports_status ON issue_reports(status);

CREATE INDEX idx_issue_reports_priority ON issue_reports(priority);

CREATE TABLE issue_report_translations (
  id VARCHAR(26) PRIMARY KEY,
  report_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  title VARCHAR(200) NOT NULL,
  description TEXT NULL,
  resolution_notes TEXT NULL,
  UNIQUE (report_id, lang_code),
  FOREIGN KEY (report_id) REFERENCES issue_reports(id) ON DELETE CASCADE
);

CREATE INDEX idx_report_translations_report_lang ON issue_report_translations(report_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_report_translations_report_lang;

DROP TABLE IF EXISTS issue_report_translations;

DROP INDEX IF EXISTS idx_issue_reports_priority;

DROP INDEX IF EXISTS idx_issue_reports_status;

DROP INDEX IF EXISTS idx_issue_reports_asset_id;

DROP TABLE IF EXISTS issue_reports;

DROP TYPE IF EXISTS issue_status;

DROP TYPE IF EXISTS issue_priority;
