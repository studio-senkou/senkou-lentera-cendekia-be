ALTER TABLE users DROP COLUMN IF EXISTS email_verified_at;
ALTER TABLE users DROP COLUMN IF EXISTS is_active;

DROP INDEX IF EXISTS idx_users_email_verified_at;
DROP INDEX IF EXISTS idx_users_is_active;
