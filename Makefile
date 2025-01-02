.PHONY: build run docker-build docker-run migrate-up migrate-down test test-local test-prod

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
	docker-compose down
	docker-compose up -d
	./test.sh
	docker-compose logs

test-prod:
	export HEROKU_APP_NAME=your-app-name && ./test.sh

migrate-up:
	docker-compose exec web goose -dir migrations postgres "host=db port=5432 user=postgres password=postgres dbname=lenslocked sslmode=disable" up

migrate-down:
	docker-compose exec web goose -dir migrations postgres "host=db port=5432 user=postgres password=postgres dbname=lenslocked sslmode=disable" down

migrate-status:
	docker-compose exec web goose -dir migrations postgres "host=db port=5432 user=postgres password=postgres dbname=lenslocked sslmode=disable" status

create-migration:
	@read -p "Enter migration name: " name; \
	docker-compose exec web goose -dir migrations create $${name} sql
