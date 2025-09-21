-- migrate:up
CREATE TABLE IF NOT EXISTS meeting_session_proofs (
    id SERIAL PRIMARY KEY,
    meeting_id INTEGER NOT NULL,
    student_proof VARCHAR(255),
    student_signature VARCHAR(255),
    mentor_proof VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

DO $$
    BEGIN

        -- Verify if the foreign key constraint for meeting_sessions already exists
        -- If it doesn't exist, continue to create new foreign key constraint
        IF NOT EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conname = 'fk_mt_session_proof'
        ) THEN
            ALTER TABLE meeting_session_proofs
            ADD CONSTRAINT fk_mt_session_proof
            FOREIGN KEY (meeting_id) REFERENCES meeting_sessions(id)
            ON UPDATE CASCADE ON DELETE CASCADE;
        END IF;

        -- Verify if the index for meeting_id already exists
        -- If it doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_meeting_session_proofs_meeting_id'
        ) THEN
            CREATE INDEX idx_meeting_session_proofs_meeting_id ON meeting_session_proofs(meeting_id);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
ALTER TABLE meeting_session_proofs
    DROP CONSTRAINT IF EXISTS fk_mt_session_proof;

DROP INDEX IF EXISTS idx_meeting_session_proofs_meeting_id;

DROP TABLE IF EXISTS meeting_session_proofs;
