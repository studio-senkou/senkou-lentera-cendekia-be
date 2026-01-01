-- migrate:up
CREATE TABLE IF NOT EXISTS blogs (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify if the author foreign key constraint exists
        -- If doens't exist, continue to create new foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_blogs_author'
        ) THEN
            ALTER TABLE blogs
            ADD CONSTRAINT fk_blogs_author
            FOREIGN KEY (author_id) REFERENCES users(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index on author_id exists
        -- If doens't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_blogs_author_id'
        ) THEN
            CREATE INDEX idx_blogs_author_id ON blogs(author_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE blogs DROP CONSTRAINT IF EXISTS fk_blogs_author;

DROP INDEX IF EXISTS idx_blogs_author_id;

DROP TABLE IF EXISTS blogs;
