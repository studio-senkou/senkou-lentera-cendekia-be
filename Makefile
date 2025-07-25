.PHONY=migrate-up migrate-down generate-app-key

include .env
export

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