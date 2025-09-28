# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/build ./cmd/server


# Stop infrastructure with Docker Compose
.PHONY: down
down:
	@docker compose down

.PHONY: up
up:
	@echo "🚀 Starting complete development environment setup..."
	@echo "\n📦 Step 1: Building Docker image..."
	@docker build -f ./Dockerfile.local -t go-backend:latest .
	@echo "\n🧹 Step 2: Cleaning old Docker images..."
	@docker image prune -f
	@echo "\n🐳 Step 3: Starting infrastructure with Docker Compose..."
	@docker compose up -d
	@echo "\n✅ All setup completed! Migrations will run automatically before starting the app."

# Clean the test cache
.PHONY: clean-testcache
clean-testcache:
	@go clean -testcache

# Show current status of all migrations
.PHONY: migrate-status
migrate-status:
	@goose -dir ./migration postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" status

# Create a new migration file with a user-provided name
.PHONY: migrate-create
migrate-create:
	@mkdir -p ./migration
	@read -p "Enter migration name: " name; \
	goose -s -dir ./migration create "$$name" sql

# Apply all up (new) migrations
.PHONY: migrate-up
migrate-up:
	@goose -dir ./migration postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" up

# Roll back the most recent migration
.PHONY: migrate-down
migrate-down:
	@goose -dir ./migration postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" down

# Generate Go code from SQL queries
.PHONY: sqlc
sqlc:
	@sqlc generate
