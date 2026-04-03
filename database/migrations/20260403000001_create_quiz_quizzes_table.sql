-- migrate:up
CREATE TABLE IF NOT EXISTS quiz_quizzes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    passing_score SMALLINT NOT NULL DEFAULT 70,   -- minimum score to pass (0-100)
    time_limit_minutes SMALLINT,                   -- NULL = no time limit
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_quizzes_is_active'
        ) THEN
            CREATE INDEX idx_quiz_quizzes_is_active ON quiz_quizzes(is_active);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
DROP INDEX IF EXISTS idx_quiz_quizzes_is_active;
DROP TABLE IF EXISTS quiz_quizzes;
