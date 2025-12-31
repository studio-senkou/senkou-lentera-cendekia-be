#!/bin/bash
set -e

USERLIST_FILE="/etc/pgbouncer/userlist.txt"

echo "Generating PgBouncer userlist..."

if [ -z "$POSTGRES_USER" ] || [ -z "$POSTGRES_PASSWORD" ]; then
    echo "Error: POSTGRES_USER and POSTGRES_PASSWORD must be set"
    exit 1
fi

PASSWORD_HASH=$(echo -n "${POSTGRES_PASSWORD}${POSTGRES_USER}" | md5sum | cut -d' ' -f1)

echo "\"${POSTGRES_USER}\" \"md5${PASSWORD_HASH}\"" > "$USERLIST_FILE"

echo "Userlist generated successfully"

exec pgbouncer /etc/pgbouncer/pgbouncer.ini
