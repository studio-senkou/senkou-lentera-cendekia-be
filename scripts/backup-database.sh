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

: "${BACKUP_RETENTION_DAYS:=7}"
: "${BACKUP_DIR:=/backups}"

TIMESTAMP=$(date +"%Y%m%d_%H%M%S")
BACKUP_FILE="${BACKUP_DIR}/${DB_DATABASE}_${TIMESTAMP}.sql.gz"
BACKUP_LATEST="${BACKUP_DIR}/${DB_DATABASE}_latest.sql.gz"

export PGPASSWORD="${DB_PASSWORD}"

echo "[$(date)] Starting backup of ${DB_DATABASE}..."

mkdir -p "${BACKUP_DIR}"

pg_dump -h "${DB_HOST}" -p "${DB_PORT}" -U "${DB_USERNAME}" -d "${DB_DATABASE}" \
    --format=plain \
    --no-owner \
    --no-privileges \
    --clean \
    --if-exists \
    | gzip > "${BACKUP_FILE}"

ln -sf "${BACKUP_FILE}" "${BACKUP_LATEST}"

echo "[$(date)] Backup completed: ${BACKUP_FILE}"

echo "[$(date)] Cleaning up old backups (older than ${BACKUP_RETENTION_DAYS} days)..."
find "${BACKUP_DIR}" -name "*.sql.gz" -type f -mtime +${BACKUP_RETENTION_DAYS} -delete

BACKUP_SIZE=$(du -h "${BACKUP_FILE}" | cut -f1)
TOTAL_BACKUPS=$(find "${BACKUP_DIR}" -name "*.sql.gz" -type f | wc -l)

echo "[$(date)] Backup size: ${BACKUP_SIZE}"
echo "[$(date)] Total backups: ${TOTAL_BACKUPS}"

if [ -n "${AWS_S3_BUCKET}" ] && [ -n "${AWS_S3_ACCESS_KEY_ID}" ]; then
    echo "[$(date)] Uploading to S3..."
    
    S3_PATH="s3://${AWS_S3_BUCKET}/backups/postgres/${DB_DATABASE}_${TIMESTAMP}.sql.gz"
    
    aws s3 cp "${BACKUP_FILE}" "${S3_PATH}" \
        --storage-class STANDARD_IA
    
    echo "[$(date)] Uploaded to: ${S3_PATH}"
    
    echo "[$(date)] Cleaning up old S3 backups..."
    aws s3 ls "s3://${AWS_S3_BUCKET}/backups/postgres/" | while read -r line; do
        FILE_DATE=$(echo "$line" | awk '{print $1}')
        FILE_NAME=$(echo "$line" | awk '{print $4}')
        
        if [ -n "$FILE_NAME" ]; then
            FILE_EPOCH=$(date -d "$FILE_DATE" +%s 2>/dev/null || echo 0)
            CUTOFF_EPOCH=$(date -d "-${BACKUP_RETENTION_DAYS} days" +%s)
            
            if [ "$FILE_EPOCH" -lt "$CUTOFF_EPOCH" ]; then
                echo "[$(date)] Deleting old S3 backup: ${FILE_NAME}"
                aws s3 rm "s3://${AWS_S3_BUCKET}/backups/postgres/${FILE_NAME}"
            fi
        fi
    done
fi

echo "[$(date)] Backup process completed successfully!"
