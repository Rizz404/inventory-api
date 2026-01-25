-- +goose Up
-- Fix constraint behaviors to follow best practices:
-- - RESTRICT for audit trails & critical references
-- - Prevent accidental data cascade deletion

-- 1. Fix assets.category_id: CASCADE -> RESTRICT
-- Asset tidak boleh terhapus otomatis saat category dihapus
ALTER TABLE assets DROP CONSTRAINT assets_category_id_fkey;

ALTER TABLE assets
ADD CONSTRAINT assets_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT;

-- 2. Fix issue_reports.reported_by: CASCADE -> RESTRICT
-- Audit trail - siapa yang melaporkan issue
ALTER TABLE issue_reports DROP CONSTRAINT issue_reports_reported_by_fkey;

ALTER TABLE issue_reports
ADD CONSTRAINT issue_reports_reported_by_fkey FOREIGN KEY (reported_by) REFERENCES users(id) ON DELETE RESTRICT;

-- 3. Fix scan_logs.scanned_by: CASCADE -> RESTRICT
-- Audit trail - siapa yang melakukan scan
ALTER TABLE scan_logs DROP CONSTRAINT scan_logs_scanned_by_fkey;

ALTER TABLE scan_logs
ADD CONSTRAINT scan_logs_scanned_by_fkey FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE RESTRICT;

-- +goose Down
-- Revert to original behavior
ALTER TABLE assets DROP CONSTRAINT assets_category_id_fkey;

ALTER TABLE assets
ADD CONSTRAINT assets_category_id_fkey FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE;

ALTER TABLE issue_reports DROP CONSTRAINT issue_reports_reported_by_fkey;

ALTER TABLE issue_reports
ADD CONSTRAINT issue_reports_reported_by_fkey FOREIGN KEY (reported_by) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE scan_logs DROP CONSTRAINT scan_logs_scanned_by_fkey;

ALTER TABLE scan_logs
ADD CONSTRAINT scan_logs_scanned_by_fkey FOREIGN KEY (scanned_by) REFERENCES users(id) ON DELETE CASCADE;
