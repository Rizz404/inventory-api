-- Add FCM token column to users table
ALTER TABLE users
ADD COLUMN fcm_token TEXT;

-- Add comment to explain the column
COMMENT ON COLUMN users.fcm_token IS 'Firebase Cloud Messaging token for push notifications';
