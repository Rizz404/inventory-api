-- +goose Up
CREATE TYPE notification_type AS ENUM (
  'MAINTENANCE',
  'WARRANTY',
  'ISSUE',
  'MOVEMENT',
  'STATUS_CHANGE',
  'LOCATION_CHANGE',
  'CATEGORY_CHANGE'
);

CREATE TYPE notification_priority AS ENUM ('LOW', 'NORMAL', 'HIGH', 'URGENT');

CREATE TABLE notifications (
  id VARCHAR(26) PRIMARY KEY,
  user_id VARCHAR(26) NOT NULL,
  -- Related entity (polymorphic approach)
  related_entity_type VARCHAR(50) NULL,
  related_entity_id VARCHAR(26) NULL,
  -- Legacy support (deprecated, use related_entity_id instead)
  related_asset_id VARCHAR(26) NULL,
  type notification_type NOT NULL,
  priority notification_priority DEFAULT 'NORMAL',
  -- Status
  is_read BOOLEAN DEFAULT FALSE,
  read_at TIMESTAMP WITH TIME ZONE NULL,
  -- Expiration
  expires_at TIMESTAMP WITH TIME ZONE NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
  FOREIGN KEY (related_asset_id) REFERENCES assets(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_notifications_user_id_is_read ON notifications(user_id, is_read);

CREATE INDEX idx_notifications_type ON notifications(type);

CREATE INDEX idx_notifications_priority ON notifications(priority);

CREATE INDEX idx_notifications_related_entity ON notifications(related_entity_type, related_entity_id);

CREATE INDEX idx_notifications_expires_at ON notifications(expires_at);

CREATE TABLE notification_translations (
  id VARCHAR(26) PRIMARY KEY,
  notification_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  title VARCHAR(200) NOT NULL,
  message TEXT NOT NULL,
  UNIQUE (notification_id, lang_code),
  FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
);

CREATE INDEX idx_notification_translations_notification_lang ON notification_translations(notification_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_notification_translations_notification_lang;

DROP TABLE IF EXISTS notification_translations;

DROP INDEX IF EXISTS idx_notifications_expires_at;

DROP INDEX IF EXISTS idx_notifications_related_entity;

DROP INDEX IF EXISTS idx_notifications_priority;

DROP INDEX IF EXISTS idx_notifications_type;

DROP INDEX IF EXISTS idx_notifications_user_id_is_read;

DROP TABLE IF EXISTS notifications;

DROP TYPE IF EXISTS notification_priority;

DROP TYPE IF EXISTS notification_type;
