-- migrate:up

-- Status attempt:
--   'in_progress' : user sudah memulai, belum submit
--   'completed'   : user sudah submit jawaban
--   'reset'       : admin mereset attempt, user bisa mulai ulang

CREATE TABLE IF NOT EXISTS quiz_attempts (
    id SERIAL PRIMARY KEY,
    quiz_id INTEGER NOT NULL,
    user_id INTEGER NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'in_progress',  -- 'in_progress', 'completed', 'reset'
    score NUMERIC(5, 2),                                 -- NULL saat in_progress, diisi saat completed
    started_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    submitted_at TIMESTAMP,                              -- NULL saat in_progress
    reset_at TIMESTAMP,                                  -- kapan admin melakukan reset
    reset_by INTEGER,                                    -- user ID admin yang me-reset
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DO $$
    BEGIN

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_attempts_quiz_id'
        ) THEN
            ALTER TABLE quiz_attempts
            ADD CONSTRAINT fk_quiz_attempts_quiz_id
            FOREIGN KEY (quiz_id) REFERENCES quiz_quizzes(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_attempts_user_id'
        ) THEN
            ALTER TABLE quiz_attempts
            ADD CONSTRAINT fk_quiz_attempts_user_id
            FOREIGN KEY (user_id) REFERENCES users(id)
            ON DELETE CASCADE;
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conname = 'fk_quiz_attempts_reset_by'
        ) THEN
            ALTER TABLE quiz_attempts
            ADD CONSTRAINT fk_quiz_attempts_reset_by
            FOREIGN KEY (reset_by) REFERENCES users(id)
            ON DELETE SET NULL;
        END IF;

        -- Composite unique: satu user hanya bisa punya 1 attempt aktif per kuis
        -- (status NOT 'reset'). Ini di-enforce di level aplikasi karena partial unique
        -- index tidak bisa secara langsung mewakili logika "satu attempt aktif".
        -- Kita buat index untuk performa query.
        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_attempts_user_quiz'
        ) THEN
            CREATE INDEX idx_quiz_attempts_user_quiz ON quiz_attempts(user_id, quiz_id);
        END IF;

        IF NOT EXISTS (
            SELECT 1 FROM pg_indexes
            WHERE indexname = 'idx_quiz_attempts_status'
        ) THEN
            CREATE INDEX idx_quiz_attempts_status ON quiz_attempts(status);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS fk_quiz_attempts_quiz_id;
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS fk_quiz_attempts_user_id;
ALTER TABLE quiz_attempts DROP CONSTRAINT IF EXISTS fk_quiz_attempts_reset_by;
DROP INDEX IF EXISTS idx_quiz_attempts_user_quiz;
DROP INDEX IF EXISTS idx_quiz_attempts_status;
DROP TABLE IF EXISTS quiz_attempts;
