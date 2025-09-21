-- migrate:up
CREATE TABLE IF NOT EXISTS mentors (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    class_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify if the foreign key constraint for users already exists
        -- If it doesn't exist, continue to create foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_mentors_user'
        ) THEN 
            ALTER TABLE mentors
            ADD CONSTRAINT fk_mentors_user
            FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the foreign key constraint for classes already exists
        -- If it doesn't exist, continue to create foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_mentors_class'
        ) THEN
            ALTER TABLE mentors
            ADD CONSTRAINT fk_mentors_class
            FOREIGN KEY (class_id) REFERENCES classes(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index for user_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_mentors_user_id'
        ) THEN
            CREATE INDEX idx_mentors_user_id ON mentors(user_id);
        END IF;

        -- Verify if the index for class_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_mentors_class_id'
        ) THEN
            CREATE INDEX idx_mentors_class_id ON mentors(class_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE mentors
    DROP CONSTRAINT IF EXISTS fk_mentors_user;

ALTER TABLE mentors
    DROP CONSTRAINT IF EXISTS fk_mentors_class;

DROP INDEX IF EXISTS idx_mentors_user_id;

DROP INDEX IF EXISTS idx_mentors_class_id;

DROP TABLE IF EXISTS mentors;