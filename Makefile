# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/app ./cmd/server

# Tidy go.mod and go.sum
.PHONY: tidy
tidy:
	@go mod tidy


# Run the application using air
.PHONY: run
run:
	@air

# Clean the test cache
.PHONY: clean-testcache
clean-testcache:
	@go clean -testcache

# Start infrastructure with Docker Compose
.PHONY: infra-up
infra-up:
	@docker compose up -d

# Stop infrastructure with Docker Compose
.PHONY: infra-down
infra-down:
	@docker compose down

.PHONY: install-goose
install-goose:
	@go install github.com/pressly/goose/v3/cmd/goose@latest

.PHONY: migrate-status
migrate-status:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" status

.PHONY: migrate-create
migrate-create:
	@mkdir -p ./migrations
	@read -p "Enter migration name: " name; \
	goose -s -dir ./migrations create "$$name" sql

.PHONY: migrate-up
migrate-up:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" up

.PHONY: migrate-down
migrate-down:
	@goose -dir ./migrations postgres "postgres://$$DB_USER:$$DB_PASSWORD@$$DB_HOST:$$DB_PORT/$$DB_NAME?sslmode=$$DB_SSL_MODE" down
