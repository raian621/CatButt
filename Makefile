VERSION := 0.0.1b
DB_CONTAINER := catbook-postgres

.PHONY: dev
dev: deps
	@echo "Starting live reload setup for CatBook v$(VERSION)"
	@air

.PHONY: build
build: deps
	@echo "Building CatButt v$(VERSION)"
	@go build

.PHONY: dev
dev: start-db
	@air

.PHONY: deps
deps:
	@go get ./...

.PHONY: db
db:
	@echo "Building docker container for postgres db..."
	@docker run --name $(DB_CONTAINER) -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=catbook_test -e POSTGRES_USER=catbook -p 5432:5432 -d postgres

.PHONY: start-db
start-db:
	@echo "Starting docker container postgres db"
	@docker start $(DB_CONTAINER)
