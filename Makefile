.PHONY: build run docker-build docker-run migrate-up migrate-down test test-local test-prod test-integration test-all deploy-prod

build:
	go build -o bin/app

run:
	go run main.go

docker-build:
	docker-compose build

docker-run:
	docker-compose up

docker-down:
	docker-compose down

heroku-create:
	heroku create
	heroku stack:set container
	heroku addons:create heroku-postgresql:mini

heroku-config:
	@echo "Setting up Heroku config vars..."
	@CSRF_KEY=$$(openssl rand -base64 32) && \
	SESSION_KEY=$$(openssl rand -base64 32) && \
	heroku config:set CSRF_KEY=$$CSRF_KEY && \
	heroku config:set SESSION_KEY=$$SESSION_KEY && \
	heroku config:set CSRF_SECURE=true && \
	heroku config:set FIVEMANAGE_DEBUG=false
	@echo "Don't forget to set your FiveManage API key:"
	@echo "heroku config:set FIVEMANAGE_API_KEY=your-key-here"

heroku-deploy:
	git push heroku main

heroku-logs:
	heroku logs --tail

heroku-db-console:
	heroku pg:psql

heroku-bash:
	heroku run bash

local-db:
	docker-compose up db

test:
	./test.sh

test-local:
	@echo "Running local environment tests..."
	@echo "Ensuring Docker services are running..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 15  # Increased wait time
	@chmod +x scripts/test/local.sh
	@./scripts/test/local.sh

test-prod:
	@echo "Running production environment tests..."
	@chmod +x scripts/test/production.sh
	@HEROKU_APP_NAME=lenslocked ./scripts/test/production.sh

test-integration:
	@echo "Running integration tests..."
	@echo "Ensuring Docker services are running..."
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 10  # Give services time to start
	@chmod +x scripts/test/integration.sh
	@./scripts/test/integration.sh

test-all: 
	@echo "Starting all tests..."
	@docker-compose down  # Ensure clean state
	@docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 15  # Increased wait time
	@make test-local
	@make test-integration
	@echo "All tests completed!"
	
# Optional: Add cleanup
test-cleanup:
	@echo "Cleaning up test environment..."
	@docker-compose down

# Add cleanup to the workflow
deploy-prod: test-all test-cleanup
	git push heroku main
	make test-prod
