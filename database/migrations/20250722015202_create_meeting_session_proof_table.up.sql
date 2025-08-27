CREATE TABLE IF NOT EXISTS meeting_session_proof (
    id SERIAL PRIMARY KEY,
    meeting_id INTEGER NOT NULL,
    student_proof VARCHAR(255),
    student_signature VARCHAR(255),
    mentor_proof VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_mt_session_proof
        FOREIGN KEY (meeting_id) REFERENCES meeting_sessions(id)
        ON UPDATE CASCADE ON DELETE CASCADE
)