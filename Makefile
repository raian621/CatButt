VERSION := 0.0.1b
DB_CONTAINER := catbook-postgres

.PHONY: dev
dev: deps dev-site dev-server
	@echo "Starting live reload setup for CatBook v$(VERSION)"

.PHONY: build
build: deps build-site build-server
	@echo "Building CatBook v$(VERSION)"

.PHONY: deps
deps: site-deps server-deps

.PHONY: dev-site
dev-site:
	@cd site; npm run dev

.PHONY: dev-server
dev-server: start-db
	@cd server; air

.PHONY: build-site
build-site:
	@cd site; npm run build

.PHONY: build-server
build-server: start-db
	@cd server; go build

.PHONY: site-deps
site-deps:
	@cd site; npm i

.PHONY: server-deps
server-deps:
	@cd server; go get ./...

.PHONY: db
db:
	@echo "Building docker container for postgres db..."
	@docker run --name $(DB_CONTAINER) -e POSTGRES_PASSWORD=postgres -e POSTGRES_DB=catbook_test -e POSTGRES_USER=catbook -p 5432:5432 -d postgres

.PHONY: start-db
start-db:
	@echo "Starting docker container postgres db"
	@docker start $(DB_CONTAINER)
