# Devon Farm Sales - Testing Makefile

.PHONY: help test test-unit test-integration test-handlers test-e2e test-all clean setup-test-db

# Default target
help: ## Show this help message
	@echo "Available targets:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  %-20s %s\n", $$1, $$2}'

# Test targets
test: test-unit ## Run unit tests (fast)

test-unit: ## Run unit tests only
	@echo "Running unit tests..."
	go test -v ./utils/... ./horse/... ./user/... ./farm/...

test-integration: ## Run integration tests (requires TEST_DATABASE_URL)
	@echo "Running integration tests..."
	@if [ -z "$(TEST_DATABASE_URL)" ]; then \
		echo "Error: TEST_DATABASE_URL environment variable not set"; \
		echo "Example: export TEST_DATABASE_URL='postgresql://root@localhost:26257/test_db?sslmode=disable'"; \
		exit 1; \
	fi
	go test -v ./tests/integration/...

test-handlers: ## Run HTTP handler tests
	@echo "Running HTTP handler tests..."
	go test -v ./tests/handlers/...

test-e2e: ## Run end-to-end tests with Playwright
	@echo "Running end-to-end tests..."
	cd tests/e2e && npm test

test-e2e-headed: ## Run end-to-end tests with visible browser
	@echo "Running end-to-end tests (headed)..."
	cd tests/e2e && npm run test:headed

test-e2e-ui: ## Open Playwright UI for test development
	cd tests/e2e && npm run test:ui

test-all: test-unit test-integration test-handlers ## Run all Go tests

# Setup targets
setup-test-db: ## Setup test database with Docker
	@echo "Starting test database..."
	docker run -d --name devon-test-db \
		-p 26258:26257 \
		-e COCKROACH_DATABASE=test_db \
		cockroachdb/cockroach:latest \
		start-single-node --insecure
	@echo "Waiting for database to be ready..."
	sleep 5
	@echo "Running migrations..."
	@export TEST_DATABASE_URL="postgresql://root@localhost:26258/test_db?sslmode=disable" && \
	export MIGRATIONS_DSN="$$TEST_DATABASE_URL" && \
	task db-migrate
	@echo "Test database ready at: postgresql://root@localhost:26258/test_db?sslmode=disable"

setup-e2e: ## Setup end-to-end test dependencies
	@echo "Installing Playwright dependencies..."
	cd tests/e2e && npm install && npx playwright install

# Cleanup targets
clean-test-db: ## Remove test database container
	@echo "Stopping and removing test database..."
	docker stop devon-test-db || true
	docker rm devon-test-db || true

clean: clean-test-db ## Clean up test resources

# Development targets
watch-unit: ## Watch for changes and run unit tests
	@echo "Watching for changes (unit tests)..."
	find . -name "*.go" -not -path "./tests/integration/*" -not -path "./tests/handlers/*" | entr -r make test-unit

watch-integration: ## Watch for changes and run integration tests
	@echo "Watching for changes (integration tests)..."
	find . -name "*.go" | entr -r make test-integration

# CI targets
ci-test: ## Run tests suitable for CI environment
	@echo "Running CI tests..."
	make test-unit
	@if [ -n "$(TEST_DATABASE_URL)" ]; then \
		make test-integration; \
	else \
		echo "Skipping integration tests (no TEST_DATABASE_URL)"; \
	fi
	make test-handlers

ci-setup: ## Setup CI environment
	@echo "Setting up CI environment..."
	# Install task runner
	go install github.com/go-task/task/v3/cmd/task@latest
	# Install go-migrate
	go install -tags 'postgres,cockroachdb' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	# Setup E2E tests
	cd tests/e2e && npm ci && npx playwright install --with-deps

# Coverage targets
coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

coverage-integration: ## Run integration tests with coverage
	@echo "Running integration tests with coverage..."
	go test -coverprofile=coverage-integration.out ./tests/integration/...
	go tool cover -html=coverage-integration.out -o coverage-integration.html

# Benchmark targets
bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

# Example usage
example-test-setup: ## Example of complete test setup
	@echo "Example: Complete test setup"
	@echo "1. Start test database:"
	@echo "   make setup-test-db"
	@echo ""
	@echo "2. Set environment variable:"
	@echo "   export TEST_DATABASE_URL='postgresql://root@localhost:26258/test_db?sslmode=disable'"
	@echo ""
	@echo "3. Run all tests:"
	@echo "   make test-all"
	@echo ""
	@echo "4. Setup E2E tests:"
	@echo "   make setup-e2e"
	@echo ""
	@echo "5. Run E2E tests:"
	@echo "   make test-e2e"
	@echo ""
	@echo "6. Clean up:"
	@echo "   make clean"