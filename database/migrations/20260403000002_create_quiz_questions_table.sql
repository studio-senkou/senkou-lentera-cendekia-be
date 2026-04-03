-- migrate:up
CREATE TABLE IF NOT EXISTS quiz_questions (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER NOT NULL,
    question_text TEXT NOT NULL,
    order_number SMALLINT NOT NULL DEFAULT 0,       -- urutan tampil soal
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_questions_quiz_id'
        ) THEN
            ALTER TABLE quiz_questions
            ADD CONSTRAINT fk_quiz_questions_quiz_id
            FOREIGN KEY (quiz_id) REFERENCES quiz_quizzes(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_questions_quiz_id'
        ) THEN
            CREATE INDEX idx_quiz_questions_quiz_id ON quiz_questions(quiz_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE quiz_questions DROP CONSTRAINT IF EXISTS fk_quiz_questions_quiz_id;
DROP INDEX IF EXISTS idx_quiz_questions_quiz_id;
DROP TABLE IF EXISTS quiz_questions;
