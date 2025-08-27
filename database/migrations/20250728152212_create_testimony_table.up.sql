CREATE TABLE IF NOT EXISTS testimonials (
    id SERIAL PRIMARY KEY,
    testimoner_name VARCHAR(255) NOT NULL,
    testimoner_current_position VARCHAR(255),
    testimoner_previous_position VARCHAR(255),
    testimoner_photo TEXT,
    testimony_text TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP,
    deleted_at TIMESTAMP
);