#!/bin/bash

set -e

if [ -f .env.production ]; then
    export $(grep -v '^#' .env.production | xargs)
fi

if [ -z "$DOMAIN" ]; then
    echo "Error: DOMAIN environment variable is not set"
    echo "Please set DOMAIN in .env.production"
    exit 1
fi

if [ -z "$CERTBOT_EMAIL" ]; then
    echo "Error: CERTBOT_EMAIL environment variable is not set"
    echo "Please set CERTBOT_EMAIL in .env.production"
    exit 1
fi

STAGING=${STAGING:-0}
RSA_KEY_SIZE=4096
DATA_PATH="./certbot"

DOMAINS=("$DOMAIN")
if [ -n "$DB_DOMAIN" ]; then
    DOMAINS+=("$DB_DOMAIN")
fi

mkdir -p "$DATA_PATH/conf"
mkdir -p "$DATA_PATH/www"

echo "### Downloading recommended TLS parameters..."
if [ ! -e "$DATA_PATH/conf/options-ssl-nginx.conf" ] || [ ! -e "$DATA_PATH/conf/ssl-dhparams.pem" ]; then
    curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot-nginx/certbot_nginx/_internal/tls_configs/options-ssl-nginx.conf > "$DATA_PATH/conf/options-ssl-nginx.conf"
    curl -s https://raw.githubusercontent.com/certbot/certbot/master/certbot/certbot/ssl-dhparams.pem > "$DATA_PATH/conf/ssl-dhparams.pem"
fi

if [ -n "$CERTBOT_EMAIL" ]; then
    EMAIL_ARG="--email $CERTBOT_EMAIL"
else
    EMAIL_ARG="--register-unsafely-without-email"
fi

if [ "$STAGING" = "1" ]; then
    STAGING_ARG="--staging"
else
    STAGING_ARG=""
fi

for CURRENT_DOMAIN in "${DOMAINS[@]}"; do
    echo ""
    echo "=========================================="
    echo "Processing domain: $CURRENT_DOMAIN"
    echo "=========================================="

    if [ -d "$DATA_PATH/conf/live/$CURRENT_DOMAIN" ]; then
        read -p "Existing certificate found for $CURRENT_DOMAIN. Replace? (y/N) " decision
        if [ "$decision" != "Y" ] && [ "$decision" != "y" ]; then
            echo "Skipping $CURRENT_DOMAIN..."
            continue
        fi
    fi

    echo "### Creating dummy certificate for $CURRENT_DOMAIN..."
    LIVE_PATH="$DATA_PATH/conf/live/$CURRENT_DOMAIN"
    mkdir -p "$LIVE_PATH"
    docker compose run --rm --entrypoint "\
      openssl req -x509 -nodes -newkey rsa:$RSA_KEY_SIZE -days 1 \
        -keyout '/etc/letsencrypt/live/$CURRENT_DOMAIN/privkey.pem' \
        -out '/etc/letsencrypt/live/$CURRENT_DOMAIN/fullchain.pem' \
        -subj '/CN=localhost'" certbot

    echo "### Starting nginx..."
    docker compose up --force-recreate -d nginx

    echo "### Deleting dummy certificate for $CURRENT_DOMAIN..."
    docker compose run --rm --entrypoint "\
      rm -Rf /etc/letsencrypt/live/$CURRENT_DOMAIN && \
      rm -Rf /etc/letsencrypt/archive/$CURRENT_DOMAIN && \
      rm -Rf /etc/letsencrypt/renewal/$CURRENT_DOMAIN.conf" certbot

    echo "### Requesting Let's Encrypt certificate for $CURRENT_DOMAIN..."
    docker compose run --rm --entrypoint "\
      certbot certonly --webroot -w /var/www/certbot \
        $STAGING_ARG \
        $EMAIL_ARG \
        -d $CURRENT_DOMAIN \
        --rsa-key-size $RSA_KEY_SIZE \
        --agree-tos \
        --force-renewal" certbot

    echo "### Certificate obtained for $CURRENT_DOMAIN!"
done

echo "### Reloading nginx..."
docker compose exec nginx nginx -s reload

echo ""
echo "=========================================="
echo "All certificates obtained successfully!"
echo "Domains: ${DOMAINS[*]}"
echo "=========================================="
