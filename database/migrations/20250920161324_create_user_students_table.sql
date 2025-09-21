-- migrate:up
CREATE TABLE IF NOT EXISTS students (
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
            WHERE conname = 'fk_students_user'
        ) THEN 
            ALTER TABLE students
            ADD CONSTRAINT fk_students_user
            FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the foreign key constraint for classes already exists
        -- If it doesn't exist, continue to create foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_students_class'
        ) THEN
            ALTER TABLE students
            ADD CONSTRAINT fk_students_class
            FOREIGN KEY (class_id) REFERENCES classes(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index already exists
        -- If it doesn't exist, continue to create index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_students_user_id'
        ) THEN
            CREATE INDEX idx_students_user_id ON students(user_id);
        END IF;

        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_students_class_id'
        ) THEN
            CREATE INDEX idx_students_class_id ON students(class_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE students
    DROP CONSTRAINT IF EXISTS fk_students_user;

ALTER TABLE students
    DROP CONSTRAINT IF EXISTS fk_students_class;

DROP INDEX IF EXISTS idx_students_user_id;

DROP INDEX IF EXISTS idx_students_class_id;

DROP TABLE IF EXISTS students;