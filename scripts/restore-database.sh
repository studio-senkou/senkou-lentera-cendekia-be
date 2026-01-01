#!/bin/bash

set -e

if [ -f .env.production ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

: "${DB_HOST:=postgres}"
: "${DB_PORT:=5432}"
: "${DB_USERNAME:=postgres}"
: "${DB_PASSWORD:?DB_PASSWORD is required}"
: "${DB_DATABASE:=senkou_lentera_cendekia}"

export PGPASSWORD="${DB_PASSWORD}"

if [ -z "$1" ]; then
    echo "Usage: $0 <backup_file.sql.gz> [--confirm]"
    echo ""
    echo "Examples:"
    echo "  $0 backups/senkou_lentera_cendekia_latest.sql.gz"
    echo "  $0 s3://bucket/backups/postgres/backup.sql.gz"
    echo ""
    echo "Options:"
    echo "  --confirm    Skip confirmation prompt"
    exit 1
fi

BACKUP_FILE="$1"
CONFIRM="$2"

if [ "$CONFIRM" != "--confirm" ]; then
    echo "WARNING: This will restore the database from: ${BACKUP_FILE}"
    echo "Database: ${DB_DATABASE} on ${DB_HOST}:${DB_PORT}"
    echo ""
    read -p "Are you sure you want to continue? (yes/no): " RESPONSE
    
    if [ "$RESPONSE" != "yes" ]; then
        echo "Restore cancelled."
        exit 0
    fi
fi

TEMP_FILE="/tmp/restore_$(date +%s).sql"

if [[ "$BACKUP_FILE" == s3://* ]]; then
    echo "[$(date)] Downloading backup from S3..."
    aws s3 cp "${BACKUP_FILE}" "${TEMP_FILE}.gz"
    gunzip "${TEMP_FILE}.gz"
else
    echo "[$(date)] Extracting backup file..."
    gunzip -c "${BACKUP_FILE}" > "${TEMP_FILE}"
fi

echo "[$(date)] Terminating existing connections to ${DB_DATABASE}..."
psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USERNAME}" -d postgres -c "
    SELECT pg_terminate_backend(pid) 
    FROM pg_stat_activity 
    WHERE datname = '${DB_DATABASE}' 
    AND pid <> pg_backend_pid();
" || true

echo "[$(date)] Restoring database ${DB_DATABASE}..."
psql -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USERNAME}" -d "${DB_DATABASE}" < "${TEMP_FILE}"

rm -f "${TEMP_FILE}"

echo "[$(date)] Restore completed successfully!"
