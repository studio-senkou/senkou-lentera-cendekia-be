.PHONY=generate-app-key generate-auth-key migrations-create migrate-up migrate-down seed rebuild-prod rebuild-dev

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

generate-auth-secret:
	@echo "Generating application auth secret.."
	@head -c 64 /dev/urandom | base64 | tr -d '\n' > auth.key
	@AUTH_SECRET=$$(cat auth.key); \
	if grep -q '^AUTH_SECRET=' .env; then \
		sed -i "s/^AUTH_SECRET=.*/AUTH_SECRET=$$AUTH_SECRET/" .env; \
	else \
		echo "AUTH_SECRET=$$AUTH_SECRET" >> .env; \
	fi

migrations-create:
	@if [ -z "$(NAME)" ]; then \
		echo "Error: Please provide a migration name using 'make migrations-create name=your_migration_name'"; \
		exit 1; \
	fi
	@dbmate -u $$(grep '^DB_URL=' .env | cut -d '=' -f2-) --migrations-dir=./database/migrations new $(NAME)

migrate-up:
	@dbmate -u $$(grep '^DB_URL=' .env | cut -d '=' -f2-) --migrations-dir=./database/migrations --schema-file=./database/migrations/schema.sql up

migrate-down:
	@dbmate -u $$(grep '^DB_URL=' .env | cut -d '=' -f2-) --migrations-dir=./database/migrations --schema-file=./database/migrations/schema.sql down

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