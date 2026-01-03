#!/bin/bash
set -e

echo "Running PostgreSQL initialization script..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    -- Enable useful extensions
    CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
    CREATE EXTENSION IF NOT EXISTS "pg_trgm";
    CREATE EXTENSION IF NOT EXISTS "btree_gin";
    
    -- Log successful initialization
    SELECT 'Database extensions initialized successfully' AS status;
EOSQL

echo "PostgreSQL initialization completed successfully."
