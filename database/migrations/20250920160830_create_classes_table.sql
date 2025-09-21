-- migrate:up
CREATE TABLE IF NOT EXISTS classes (
    id UUID PRIMARY KEY,
    classname VARCHAR(100) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- migrate:down
DROP TABLE IF EXISTS classes;
