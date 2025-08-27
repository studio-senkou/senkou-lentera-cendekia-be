CREATE TABLE IF NOT EXISTS student_plans (
    id SERIAL PRIMARY KEY,
    student_id INTEGER NOT NULL,
    total_sessions INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_student_plans_student
        FOREIGN KEY (student_id) REFERENCES students(id) ON DELETE CASCADE
)