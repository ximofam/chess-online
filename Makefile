run-api:
	go run cmd/api/main.go
.PHONY: run-api

build-and-run:
	docker-compose up -d --build
.PHONY: build-and-run

run:
	docker-compose up -d

run-docker-dev:
	docker-compose -f docker-compose.dev.yml up -d
.PHONY: run-docker-dev

