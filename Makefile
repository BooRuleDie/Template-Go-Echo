# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/build ./cmd/server

.PHONY: up
up:
	@docker build -f ./docker/Dockerfile.local -t echo_template:local .
	@docker compose -p echo_template_local -f ./docker/docker-compose-local.yml down
	@docker compose -p echo_template_local -f ./docker/docker-compose-local.yml up -d
	@docker image prune -f

.PHONY: down
down:
	@docker compose -p echo_template_local -f ./docker/docker-compose-local.yml down

.PHONY: up-dev
up-dev:
	@docker build -f ./docker/Dockerfile.dev -t echo_template:dev .
	@docker compose -p echo_template_dev -f ./docker/docker-compose-dev.yml down
	@docker compose -p echo_template_dev -f ./docker/docker-compose-dev.yml up -d
	@docker image prune -f

.PHONY: down-dev
down-dev:
	@docker compose -p echo_template_dev -f ./docker/docker-compose-dev.yml down

.PHONY: up-prod
up-prod:
	@docker build -f ./docker/Dockerfile.prod -t echo_template:prod .
	@docker compose -p echo_template_prod -f ./docker/docker-compose-prod.yml down
	@docker compose -p echo_template_prod -f ./docker/docker-compose-prod.yml up -d
	@docker image prune -f

.PHONY: down-prod
down-prod:
	@docker compose -p echo_template_prod -f ./docker/docker-compose-prod.yml down

# Clean the test cache
.PHONY: clean-testcache
clean-testcache:
	@go clean -testcache

# Show current status of all migrations
.PHONY: migrate-status
migrate-status:
	@goose -dir ./migration postgres "$$DATABASE_URL" status

# Create a new migration file with a user-provided name
.PHONY: migrate-create
migrate-create:
	@mkdir -p ./migration
	@read -p "Enter migration name: " name; \
	goose -s -dir ./migration create "$$name" sql

# Apply all up (new) migrations
.PHONY: migrate-up
migrate-up:
	@goose -dir ./migration postgres "$$DATABASE_URL" up

# Roll back the most recent migration
.PHONY: migrate-down
migrate-down:
	@goose -dir ./migration postgres "$$DATABASE_URL" down

# Generate Go code from SQL queries
.PHONY: sqlc
sqlc:
	@sqlc generate
