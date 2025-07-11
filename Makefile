.PHONY=

generate-app-key:
	@echo "Generating application key.."
	@head -c 32 /dev/urandom | base64 | tr -d '\n' > app.key
	@APP_KEY=$$(cat app.key); \
	if grep -q '^APP_KEY=' .env; then \
		sed -i "s/^APP_KEY=.*/APP_KEY=$$APP_KEY/" .env; \
	else \
		echo "APP_KEY=$$APP_KEY" >> .env; \
	fi