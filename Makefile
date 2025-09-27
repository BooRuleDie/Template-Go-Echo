# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/app ./cmd/server


# Stop infrastructure with Docker Compose
.PHONY: infra-down
infra-down:
	@docker compose down

.PHONY: run
run:
	@echo "🚀 Starting complete development environment setup..."
	@echo "\n📦 Step 1: Building Docker image..."
	@docker build -f ./Dockerfile.local -t go-backend:latest .
	@echo "\n🧹 Step 2: Cleaning old Docker images..."
	@docker image prune -f
	@echo "🐳 Step 3: Starting infrastructure with Docker Compose..."
	@docker compose up -d
	@until docker exec postgres pg_isready -U $(DB_USER) -d $(DB_NAME) >/dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@until docker exec redis redis-cli -a $(REDIS_PASSWORD) ping >/dev/null 2>&1; do \
		printf "."; \
		sleep 1; \
	done
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" up
	@echo "\n✅ All setup completed! Your development environment is ready."

# Clean the test cache
.PHONY: clean-testcache
clean-testcache:
	@go clean -testcache

# Show current status of all migrations
.PHONY: migrate-status
migrate-status:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" status

# Create a new migration file with a user-provided name
.PHONY: migrate-create
migrate-create:
	@mkdir -p ./migrations
	@read -p "Enter migration name: " name; \
	goose -s -dir ./migrations create "$$name" sql

# Apply all up (new) migrations
.PHONY: migrate-up
migrate-up:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" up

# Roll back the most recent migration
.PHONY: migrate-down
migrate-down:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" down

# Install 'sqlc' code generation tool
.PHONY: install-sqlc
install-sqlc:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

# Generate Go code from SQL queries
.PHONY: sqlc
sqlc:
	@sqlc generate
