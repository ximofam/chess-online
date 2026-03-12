run:
	docker-compose up -d
.PHONY: run

build-and-run:
	docker-compose up -d --build
.PHONY: build-and-run

run-dev:
	docker-compose -f docker-compose.dev.yml up -d

build-and-run-dev:
	docker-compose -f docker-compose.dev.yml up -d --build
.PHONY: build-and-run-dev

