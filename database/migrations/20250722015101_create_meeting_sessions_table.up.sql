CREATE TABLE IF NOT EXISTS meeting_sessions (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    mentor_id INTEGER NOT NULL,
    session_date DATE NOT NULL,
    session_time TIME NOT NULL,
    duration_minutes SMALLINT NOT NULL, -- Duration in minutes
    status VARCHAR(20) NOT NULL, -- e.g., scheduled, completed, missed, canceled
    note TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_student_mt_session
        FOREIGN KEY (student_id) REFERENCES students(id)
        ON UPDATE CASCADE ON DELETE CASCADE,
    CONSTRAINT fk_mentor_mt_session
        FOREIGN KEY (mentor_id) REFERENCES mentors(id)
        ON UPDATE CASCADE ON DELETE CASCADE
)