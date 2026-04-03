-- migrate:up
CREATE TABLE IF NOT EXISTS quiz_options (
    id SERIAL PRIMARY KEY,
    question_id INTEGER NOT NULL,
    option_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    order_number SMALLINT NOT NULL DEFAULT 0,       -- urutan tampil pilihan
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_options_question_id'
        ) THEN
            ALTER TABLE quiz_options
            ADD CONSTRAINT fk_quiz_options_question_id
            FOREIGN KEY (question_id) REFERENCES quiz_questions(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_options_question_id'
        ) THEN
            CREATE INDEX idx_quiz_options_question_id ON quiz_options(question_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE quiz_options DROP CONSTRAINT IF EXISTS fk_quiz_options_question_id;
DROP INDEX IF EXISTS idx_quiz_options_question_id;
DROP TABLE IF EXISTS quiz_options;
