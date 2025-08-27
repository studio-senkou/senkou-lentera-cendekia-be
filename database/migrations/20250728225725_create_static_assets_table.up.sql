CREATE TABLE IF NOT EXISTS static_assets (
    id SERIAL PRIMARY KEY,
    asset_name VARCHAR(255) NOT NULL,
    asset_type VARCHAR(50) NOT NULL, -- e.g., 'image', 'video', 'document'
    asset_url TEXT NOT NULL,
    asset_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);