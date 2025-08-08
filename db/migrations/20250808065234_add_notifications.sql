-- +goose Up
CREATE TYPE notification_type AS ENUM (
  'MAINTENANCE',
  'WARRANTY',
  'STATUS_CHANGE',
  'MOVEMENT',
  'ISSUE_REPORT'
);

CREATE TABLE notifications (
  id VARCHAR(26) PRIMARY KEY,
  user_id VARCHAR(26) NOT NULL,
  related_asset_id VARCHAR(26) NULL,
  type notification_type NOT NULL,
  is_read BOOLEAN DEFAULT FALSE,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (related_asset_id) REFERENCES assets(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);

CREATE INDEX idx_notifications_type ON notifications(type);

CREATE TABLE notifications_translation (
  id VARCHAR(26) PRIMARY KEY,
  notification_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  title VARCHAR(200) NOT NULL,
  message TEXT NOT NULL,
  UNIQUE (notification_id, lang_code),
  FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
);

CREATE INDEX idx_notifications_translation_notification_lang ON notifications_translation(notification_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_notifications_translation_notification_lang;

DROP TABLE IF EXISTS notifications_translation;

DROP INDEX IF EXISTS idx_notifications_type;

DROP INDEX IF EXISTS idx_notifications_user_id_is_read;

DROP TABLE IF EXISTS notifications;

DROP TYPE IF EXISTS notification_type;
