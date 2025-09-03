# Run tests only for packages with test files
.PHONY: unit-test
unit-test:
	@go test -v $(shell go list -f '{{if .TestGoFiles}}{{.ImportPath}}{{end}}' ./...)

# Build the application
.PHONY: build
build:
	@go build -o tmp/app ./cmd/server

# Run the application using air
.PHONY: run
run:
	@air

# Clean the test cache
.PHONY: clean-testcache
clean-testcache:
	@go clean -testcache
