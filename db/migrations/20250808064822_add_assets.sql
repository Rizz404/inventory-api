-- +goose Up
CREATE TYPE asset_status AS ENUM ('Active', 'Maintenance', 'Disposed', 'Lost');

CREATE TYPE asset_condition AS ENUM ('Good', 'Fair', 'Poor', 'Damaged');

CREATE TABLE assets (
  id VARCHAR(26) PRIMARY KEY,
  asset_tag VARCHAR(50) UNIQUE NOT NULL,
  data_matrix_image_url VARCHAR(255) NULL,
  asset_name VARCHAR(200) NOT NULL,
  category_id VARCHAR(26) NOT NULL,
  brand VARCHAR(100) NULL,
  model VARCHAR(100) NULL,
  serial_number VARCHAR(100) UNIQUE NULL,
  purchase_date DATE NULL,
  purchase_price DECIMAL(15, 2) NULL,
  vendor_name VARCHAR(150) NULL,
  warranty_end DATE NULL,
  status asset_status DEFAULT 'Active',
  condition_status asset_condition DEFAULT 'Good',
  location_id VARCHAR(26) NULL,
  assigned_to VARCHAR(26) NULL,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY (category_id) REFERENCES categories(id) ON DELETE RESTRICT,
  FOREIGN KEY (location_id) REFERENCES locations(id) ON DELETE
  SET NULL,
    FOREIGN KEY (assigned_to) REFERENCES users(id) ON DELETE
  SET NULL
);

CREATE INDEX idx_assets_status ON assets(status);

CREATE INDEX idx_assets_location ON assets(location_id);

CREATE INDEX idx_assets_assigned_to ON assets(assigned_to);

CREATE INDEX idx_assets_category_id ON assets(category_id);

CREATE INDEX idx_assets_warranty_end ON assets(warranty_end);

CREATE INDEX idx_assets_name_brand_model ON assets(asset_name, brand, model);

-- +goose Down
DROP INDEX IF EXISTS idx_assets_name_brand_model;

DROP INDEX IF EXISTS idx_assets_warranty_end;

DROP INDEX IF EXISTS idx_assets_category_id;

DROP INDEX IF EXISTS idx_assets_assigned_to;

DROP INDEX IF EXISTS idx_assets_location;

DROP INDEX IF EXISTS idx_assets_status;

DROP TABLE IF EXISTS assets;

DROP TYPE IF EXISTS asset_condition;

DROP TYPE IF EXISTS asset_status;
