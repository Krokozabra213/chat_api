POSTGRES_DSN=postgres://myuser:mypassword@localhost:5555/postgres?sslmode=disable
MIGRATIONS_GOOSE_DIR=migrations/goose

.PHONY: up down help migrate-create migrate-up migrate-down migrate-status migrate-reset tests

# Default target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  docker-up                            Start all containers"
	@echo "  docker-down                          Stop all containers"
	@echo "  migrate-create name=<table_name>     Create new migration file"
	@echo "  migrate-up                           Apply all pending migrations"
	@echo "  migrate-down                         Rollback last migration"
	@echo "  migrate-status                       Show migrations status"
	@echo "  migrate-reset                        Rollback all migrations"
	@echo "  tests                                Start tests"

# Start containers
docker-up:
	docker compose up -d

# Stop containers
docker-down:
	docker compose down


# Create new migration file: make migrate-create name=name_table
migrate-create:
	goose -dir $(MIGRATIONS_GOOSE_DIR) create $(name) sql

migrate-up:
	goose -dir $(MIGRATIONS_GOOSE_DIR) postgres "$(POSTGRES_DSN)" up

migrate-down:
	goose -dir $(MIGRATIONS_GOOSE_DIR) postgres "$(POSTGRES_DSN)" down

migrate-status:
	goose -dir $(MIGRATIONS_GOOSE_DIR) postgres "$(POSTGRES_DSN)" status

migrate-reset:
	goose -dir $(MIGRATIONS_GOOSE_DIR) postgres "$(POSTGRES_DSN)" reset

# Start tests
tests:
	go test -v -count=1 ./tests/...  
