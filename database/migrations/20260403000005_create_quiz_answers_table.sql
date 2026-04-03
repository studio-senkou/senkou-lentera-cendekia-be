-- migrate:up
CREATE TABLE IF NOT EXISTS quiz_answers (
    id SERIAL PRIMARY KEY,
    attempt_id INTEGER NOT NULL,
    question_id INTEGER NOT NULL,
    option_id INTEGER NOT NULL,             -- pilihan jawaban user
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_answers_attempt_id'
        ) THEN
            ALTER TABLE quiz_answers
            ADD CONSTRAINT fk_quiz_answers_attempt_id
            FOREIGN KEY (attempt_id) REFERENCES quiz_attempts(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_answers_question_id'
        ) THEN
            ALTER TABLE quiz_answers
            ADD CONSTRAINT fk_quiz_answers_question_id
            FOREIGN KEY (question_id) REFERENCES quiz_questions(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_answers_option_id'
        ) THEN
            ALTER TABLE quiz_answers
            ADD CONSTRAINT fk_quiz_answers_option_id
            FOREIGN KEY (option_id) REFERENCES quiz_options(id)
            ON DELETE CASCADE;
        END IF;

        -- Satu attempt hanya boleh punya satu jawaban per pertanyaan
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'uq_quiz_answers_attempt_question'
        ) THEN
            ALTER TABLE quiz_answers
            ADD CONSTRAINT uq_quiz_answers_attempt_question
            UNIQUE (attempt_id, question_id);
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_answers_attempt_id'
        ) THEN
            CREATE INDEX idx_quiz_answers_attempt_id ON quiz_answers(attempt_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE quiz_answers DROP CONSTRAINT IF EXISTS fk_quiz_answers_attempt_id;
ALTER TABLE quiz_answers DROP CONSTRAINT IF EXISTS fk_quiz_answers_question_id;
ALTER TABLE quiz_answers DROP CONSTRAINT IF EXISTS fk_quiz_answers_option_id;
ALTER TABLE quiz_answers DROP CONSTRAINT IF EXISTS uq_quiz_answers_attempt_question;
DROP INDEX IF EXISTS idx_quiz_answers_attempt_id;
DROP TABLE IF EXISTS quiz_answers;
