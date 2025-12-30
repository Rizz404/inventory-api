-- +goose Up
ALTER TABLE users
ADD COLUMN phone_number VARCHAR(20) NULL;

ALTER TABLE users
ADD COLUMN last_login TIMESTAMP WITH TIME ZONE NULL;

CREATE INDEX idx_users_phone_number ON users(phone_number);

-- +goose Down
DROP INDEX IF EXISTS idx_users_phone_number;

ALTER TABLE users DROP COLUMN IF EXISTS last_login;

ALTER TABLE users DROP COLUMN IF EXISTS phone_number;
