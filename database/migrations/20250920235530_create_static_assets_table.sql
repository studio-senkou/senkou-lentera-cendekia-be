-- migrate:up
CREATE TABLE IF NOT EXISTS static_assets (
    id SERIAL PRIMARY KEY,
    asset_name VARCHAR(255) NOT NULL,
    asset_type VARCHAR(50) NOT NULL, -- e.g., 'image', 'video', 'document'
    asset_url TEXT NOT NULL,
    asset_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL
);

DO $$
    BEGIN

        -- Verify if the index on asset_url exists
        -- If doesn't exist, continue to create new index
        IF NOT EXISTS (
            SELECT 1
            FROM pg_indexes
            WHERE indexname = 'idx_static_assets_url'
        ) THEN
            CREATE INDEX idx_static_assets_url ON static_assets(asset_url);
        END IF;

    END;
$$ LANGUAGE plpgsql;

-- migrate:down
DROP TABLE IF EXISTS static_assets;
