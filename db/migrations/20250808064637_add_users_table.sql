-- +goose Up
CREATE TYPE user_role AS ENUM ('Admin', 'Staff', 'Employee');

CREATE TABLE users (
  id VARCHAR(26) PRIMARY KEY,
  name VARCHAR(50) UNIQUE NOT NULL,
  email VARCHAR(255) UNIQUE NOT NULL,
  password_hash VARCHAR(255) NOT NULL,
  full_name VARCHAR(100) NOT NULL,
  role user_role NOT NULL,
  employee_id VARCHAR(20) UNIQUE NULL,
  preferred_lang VARCHAR(5) DEFAULT 'id-ID',
  is_active BOOLEAN DEFAULT TRUE,
  avatar_url VARCHAR(255) NULL,
  fcm_token TEXT NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_role_active ON users(role, is_active);

-- +goose Down
DROP INDEX IF EXISTS idx_users_role_active;

DROP TABLE IF EXISTS users;

DROP TYPE IF EXISTS user_role;
