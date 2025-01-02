.PHONY: build run docker-build docker-run migrate-up migrate-down

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

heroku-deploy:
	git push heroku main

local-db:
	docker-compose up db
