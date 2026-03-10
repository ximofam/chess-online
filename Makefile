run-api:
	go run cmd/api/main.go
.PHONY: run-api

run:
	docker-compose up -d --build
.PHONY: run

run-docker-dev:
	docker-compose -f docker-compose.dev.yml up -d
.PHONY: run-docker-dev

