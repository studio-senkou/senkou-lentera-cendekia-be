ALTER TABLE users ADD COLUMN email_verified_at TIMESTAMP DEFAULT NULL;
ALTER TABLE users ADD COLUMN is_active BOOLEAN DEFAULT FALSE;

CREATE INDEX idx_users_email_verified_at ON users(email_verified_at);
CREATE INDEX idx_users_is_active ON users(is_active);

COMMENT ON COLUMN users.email_verified_at IS 'Timestamp when user email was verified';
COMMENT ON COLUMN users.is_active IS 'Whether user account is active';