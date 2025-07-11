.PHONY=migrate-up generate-app-key

generate-app-key:
	@echo "Generating application key.."
	@head -c 32 /dev/urandom | base64 | tr -d '\n' > app.key
	@APP_KEY=$$(cat app.key); \
	if grep -q '^APP_KEY=' .env; then \
		sed -i "s/^APP_KEY=.*/APP_KEY=$$APP_KEY/" .env; \
	else \
		echo "APP_KEY=$$APP_KEY" >> .env; \
	fi

DB_URL=postgres://$(DB_USERNAME):$(shell printf '%s' "$(DB_PASSWORD)" | sed -e 's/!/%21/g' -e 's/#/%23/g' -e 's/@/%40/g' -e 's/\$$/%24/g')@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable

migrate-up:
	@echo "Migrating database up.."
	@migrate -path database/migrations -database "$(DB_URL)" up

migrate-down:
	@echo "Migrating database down.."
	@migrate -path database/migrations -database "$(DB_URL)" down