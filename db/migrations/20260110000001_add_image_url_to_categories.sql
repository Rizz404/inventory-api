-- +goose Up
-- +goose StatementBegin
ALTER TABLE categories
ADD COLUMN image_url TEXT;

COMMENT ON COLUMN categories.image_url IS 'URL of category image stored in Cloudinary';

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE categories DROP COLUMN image_url;

-- +goose StatementEnd
