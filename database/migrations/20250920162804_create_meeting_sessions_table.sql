-- migrate:up
CREATE TABLE IF NOT EXISTS meeting_sessions (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    mentor_id INTEGER NOT NULL,
    session_date DATE NOT NULL,
    session_time TIME NOT NULL,
    duration_minutes SMALLINT NOT NULL, -- Duration in minutes
    status VARCHAR(20) NOT NULL, -- e.g., scheduled, completed, missed, canceled
    note TEXT,
    description TEXT,
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
            WHERE conname = 'fk_student_mt_sessions'
        ) THEN
            ALTER TABLE meeting_sessions
            ADD CONSTRAINT fk_student_mt_sessions
            FOREIGN KEY (student_id) REFERENCES students(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the foreign key constraint for mentors already exists
        -- If it doesn't exist, continue to create new foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_mentor_mt_sessions'
        ) THEN
            ALTER TABLE meeting_sessions
            ADD CONSTRAINT fk_mentor_mt_sessions
            FOREIGN KEY (mentor_id) REFERENCES mentors(id)
            ON DELETE CASCADE;
        END IF;

        -- Verify if the index for student_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_meeting_sessions_student_id'
        ) THEN
            CREATE INDEX idx_meeting_sessions_student_id ON meeting_sessions(student_id);
        END IF;

        -- Verify if the index for mentor_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_meeting_sessions_mentor_id'
        ) THEN
            CREATE INDEX idx_meeting_sessions_mentor_id ON meeting_sessions(mentor_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE meeting_sessions
    DROP CONSTRAINT IF EXISTS fk_student_mt_sessions;

ALTER TABLE meeting_sessions
    DROP CONSTRAINT IF EXISTS fk_mentor_mt_sessions;

DROP INDEX IF EXISTS idx_meeting_sessions_student_id;

DROP INDEX IF EXISTS idx_meeting_sessions_mentor_id;

DROP TABLE IF EXISTS meeting_sessions;