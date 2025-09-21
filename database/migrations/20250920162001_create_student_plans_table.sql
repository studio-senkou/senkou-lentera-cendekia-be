-- migrate:up
CREATE TABLE IF NOT EXISTS student_plans (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    total_sessions INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify if the foreign key constraint for students already exists
        -- If it doesn't exist, continue to create new foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_student_plans'
        ) THEN
            ALTER TABLE student_plans
            ADD CONSTRAINT fk_student_plans
            FOREIGN KEY (student_id) REFERENCES students(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index for student_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_student_plans_student_id'
        ) THEN
            CREATE INDEX idx_student_plans_student_id ON student_plans(student_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE student_plans
    DROP CONSTRAINT IF EXISTS fk_student_plans;

DROP INDEX IF EXISTS idx_student_plans_student_id;

DROP TABLE IF EXISTS student_plans;