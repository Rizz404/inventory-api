-- +goose Up
CREATE TABLE categories (
  id VARCHAR(26) PRIMARY KEY,
  parent_id VARCHAR(26) NULL,
  category_code VARCHAR(20) UNIQUE NOT NULL,
  CONSTRAINT fk_parent_category FOREIGN KEY (parent_id) REFERENCES categories(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_categories_parent_id ON categories(parent_id);

CREATE TABLE categories_translation (
  id VARCHAR(26) PRIMARY KEY,
  category_id VARCHAR(26) NOT NULL,
  lang_code VARCHAR(5) NOT NULL,
  category_name VARCHAR(100) NOT NULL,
  description TEXT NULL,
  UNIQUE (category_id, lang_code),
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE CASCADE
);

CREATE INDEX idx_categories_translation_category_lang ON categories_translation(category_id, lang_code);

-- +goose Down
DROP INDEX IF EXISTS idx_categories_translation_category_lang;

DROP TABLE IF EXISTS categories_translation;

DROP INDEX IF EXISTS idx_categories_parent_id;

DROP TABLE IF EXISTS categories;
