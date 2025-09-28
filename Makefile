# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/build ./cmd/server

.PHONY: docker-network
docker-network:
	@docker network create echo_template

.PHONY: infra-up
infra-up:
	@docker compose -f docker-compose.infra.yml up -d

.PHONY: infra-up-local
infra-up-local:
	@ENV_TYPE=local $(MAKE) --no-print-directory infra-up

.PHONY: infra-up-dev
infra-up-dev:
	@ENV_TYPE=dev $(MAKE) --no-print-directory infra-up

.PHONY: infra-up-prod
infra-up-prod:
	@ENV_TYPE=prod $(MAKE) --no-print-directory infra-up

.PHONY: infra-down
infra-down:
	@docker compose -f docker-compose.infra.yml down

.PHONY: infra-down-local
infra-down-local:
	@ENV_TYPE=local $(MAKE) --no-print-directory infra-down

.PHONY: infra-down-dev
infra-down-dev:
	@ENV_TYPE=dev $(MAKE) --no-print-directory infra-down

.PHONY: infra-down-prod
infra-down-prod:
	@ENV_TYPE=prod $(MAKE) --no-print-directory infra-down

.PHONY: up
up:
	@echo "🚀 Starting complete development environment setup..."
	@echo "\n🐳 Step 1: Starting infrastructure containers..."
	@$(MAKE) --no-print-directory infra-up-local
	@sleep 3
	@echo "\n📦 Step 2: Building Docker image..."
	@docker build -f ./Dockerfile.local -t echo-template:local .
	@echo "\n🧹 Step 3: Cleaning old Docker images..."
	@docker image prune -f
	@echo "\n🐳 Step 4: Starting application containers with Docker Compose..."
	@ENV_FILE=.env.local ENV_TYPE=local docker compose -f docker-compose.local.yml up -d
	@echo "\n✅ All setup completed! Migrations will run automatically before starting the app."

.PHONY: down
down:
	@ENV_FILE=.env.local ENV_TYPE=local docker compose -f docker-compose.local.yml down
	@$(MAKE) --no-print-directory infra-down-local

.PHONY: up-dev
up-dev:
	@echo "🚀 Starting development environment setup..."
	@echo "\n🐳 Step 1: Starting development infrastructure containers..."
	@$(MAKE) --no-print-directory infra-up-dev
	@sleep 3
	@echo "\n📦 Step 2: Building Docker image for development..."
	@ENV_FILE=.env.dev ENV_TYPE=dev docker build -t echo-template:dev .
	@echo "\n🧹 Step 3: Cleaning old Docker images..."
	@docker image prune -f
	@echo "\n🐳 Step 4: Starting development containers with Docker Compose..."
	@ENV_FILE=.env.dev ENV_TYPE=dev docker compose up -d
	@echo "\n✅ Development environment is ready!"

.PHONY: down-dev
down-dev:
	@ENV_FILE=.env.dev ENV_TYPE=dev docker compose down
	@$(MAKE) --no-print-directory infra-down-dev

.PHONY: up-prod
up-prod:
	@echo "🚀 Starting production environment setup..."
	@echo "\n🐳 Step 1: Starting production infrastructure containers..."
	@$(MAKE) --no-print-directory infra-up-prod
	@sleep 3
	@echo "\n📦 Step 2: Building optimized Docker image for production..."
	@ENV_FILE=.env.prod ENV_TYPE=prod docker build -t echo-template:prod .
	@echo "\n🧹 Step 3: Cleaning old Docker images..."
	@docker image prune -f
	@echo "\n🐳 Step 4: Starting production containers with Docker Compose..."
	@ENV_FILE=.env.prod ENV_TYPE=prod docker compose up -d
	@echo "\n✅ Production environment is ready!"

# Stop production infrastructure
.PHONY: down-prod
down-prod:
	@ENV_FILE=.env.prod ENV_TYPE=prod docker compose down
	@$(MAKE) --no-print-directory infra-down-prod

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
