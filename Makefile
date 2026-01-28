.PHONY: up down help

# Default target
help:
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  up        Start all containers"
	@echo "  down      Stop all containers"

# Start containers
up:
	docker compose up -d

# Stop containers
down:
	docker compose down
