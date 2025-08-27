CREATE TABLE IF NOT EXISTS mentors (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    class_id UUID NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP,

    CONSTRAINT fk_mentors_user
        FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_mentors_class
        FOREIGN KEY (class_id) REFERENCES classes(id) ON DELETE CASCADE
);