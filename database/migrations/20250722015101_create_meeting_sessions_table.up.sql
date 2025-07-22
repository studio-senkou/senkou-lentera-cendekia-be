CREATE TABLE meeting_sessions (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    mentor_id INTEGER NOT NULL,
    session_date TIMESTAMP NOT NULL,
    session_time TIME NOT NULL,
    session_duration INT NOT NULL, -- Duration in minutes
    session_type VARCHAR(50) NOT NULL,
    session_topic VARCHAR(255) NOT NULL,
    session_description TEXT,
    session_proof TEXT,
    session_feedback TEXT,
    student_attendance_proof TEXT,
    mentor_attendance_proof TEXT,
    session_status VARCHAR(50) NOT NULL DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
)