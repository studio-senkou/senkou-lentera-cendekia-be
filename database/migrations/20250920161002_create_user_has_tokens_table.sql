-- migrate:up
CREATE TABLE IF NOT EXISTS user_has_tokens (
    id SERIAL PRIMARY KEY,
    user_id SERIAL NOT NULL UNIQUE,
    token TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify if the foreign key constraint already exists to avoid duplication
        -- If it doesn't exist, continue to create foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_user_tokens'
        ) THEN 
            ALTER TABLE user_has_tokens
            ADD CONSTRAINT fk_user_tokens
            FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index already exists
        -- If it doesn't exist, continue to create index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_user_tokens'
        ) THEN
            CREATE INDEX idx_user_tokens ON user_has_tokens(user_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE user_has_tokens
    DROP CONSTRAINT IF EXISTS fk_user_tokens;

DROP INDEX IF EXISTS idx_user_tokens;

DROP TABLE IF EXISTS user_has_tokens;