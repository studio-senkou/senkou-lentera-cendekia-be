.PHONY=migrate-up migrate-up-force migrate-down generate-app-key create-migration migrate-fresh seed rebuild-prod rebuild-dev

# Comment if want to rebuild docker containers to remove collision with environment variables
ifneq (,$(wildcard ./.env))
DB_USERNAME := $(shell grep '^DB_USERNAME=' .env | cut -d '=' -f2-)
DB_PASSWORD := $(shell grep '^DB_PASSWORD=' .env | cut -d '=' -f2-)
DB_HOST := $(shell grep '^DB_HOST=' .env | cut -d '=' -f2-)
DB_PORT := $(shell grep '^DB_PORT=' .env | cut -d '=' -f2-)
DB_DATABASE := $(shell grep '^DB_DATABASE=' .env | cut -d '=' -f2-)
export DB_USERNAME
export DB_PASSWORD
export DB_HOST
export DB_PORT
export DB_DATABASE
endif

generate-app-key:
	@echo "Generating application key.."
	@head -c 32 /dev/urandom | base64 | tr -d '\n' > app.key
	@APP_KEY=$$(cat app.key); \
	if grep -q '^APP_KEY=' .env; then \
		sed -i "s/^APP_KEY=.*/APP_KEY=$$APP_KEY/" .env; \
	else \
		echo "APP_KEY=$$APP_KEY" >> .env; \
	fi

migrate-up:
	@echo "Migrating database up.."
	@DB_PASSWORD_ENCODED=$$(printf '%s' '$(DB_PASSWORD)' | sed -e 's/!/%21/g' -e 's/#/%23/g' -e 's/@/%40/g' -e 's/\$$/%24/g'); \
	DB_URL="postgres://$(DB_USERNAME):$$DB_PASSWORD_ENCODED@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"; \
	migrate -path database/migrations -database "$$DB_URL" up

migrate-down:
	@echo "Migrating database down.."
	@DB_PASSWORD_ENCODED=$$(printf '%s' '$(DB_PASSWORD)' | sed -e 's/!/%21/g' -e 's/#/%23/g' -e 's/@/%40/g' -e 's/\$$/%24/g'); \
	DB_URL="postgres://$(DB_USERNAME):$$DB_PASSWORD_ENCODED@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"; \
	migrate -path database/migrations -database "$$DB_URL" down

create-migration:
	@echo "Creating new migration file.."
	@read -p "Enter migration name: " MIGRATION_NAME; \
	if [ -z "$$MIGRATION_NAME" ]; then \
		echo "Migration name cannot be empty"; \
		exit 1; \
	fi; \
	migrate create -ext sql -dir database/migrations "$$MIGRATION_NAME"

migrate-fresh:
	@echo "Running fresh migration.."
	@DB_PASSWORD_ENCODED=$$(printf '%s' '$(DB_PASSWORD)' | sed -e 's/!/%21/g' -e 's/#/%23/g' -e 's/@/%40/g' -e 's/\$$/%24/g'); \
	DB_URL="postgres://$(DB_USERNAME):$$DB_PASSWORD_ENCODED@$(DB_HOST):$(DB_PORT)/$(DB_DATABASE)?sslmode=disable"; \
	migrate -path database/migrations -database "$$DB_URL" down -all; \
	migrate -path database/migrations -database "$$DB_URL" up

seed:
	@echo "Seeding database.."
	@go run database/seeders/database_seeder.go
	@echo "Database seeding completed."

rebuild-prod:
	@mv .env .env.temp
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-lentera-cendekia-api down --remove-orphans
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-lentera-cendekia-api build --no-cache
	@docker compose -f docker-compose.yml --env-file .env.production -p senkou-lentera-cendekia-api up -d --force-recreate
	@mv .env.temp .env
	@echo "Rebuild completed."
	
rebuild-dev:
	@docker compose -f docker-compose.dev.yml -p senkou-lentera-cendekia-api-dev down --remove-orphans
	@docker compose -f docker-compose.dev.yml -p senkou-lentera-cendekia-api-dev build --no-cache
	@docker compose -f docker-compose.dev.yml -p senkou-lentera-cendekia-api-dev up -d --force-recreate
	@echo "Rebuild for development completed."