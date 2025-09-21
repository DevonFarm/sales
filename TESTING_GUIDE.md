# Testing Guide for Devon Farm Sales

This guide provides comprehensive instructions for testing the Devon Farm Sales application.

## Quick Start

1. **Run Unit Tests** (fastest):
   ```bash
   make test
   ```

2. **Setup Test Database and Run All Tests**:
   ```bash
   make setup-test-db
   export TEST_DATABASE_URL="postgresql://root@localhost:26258/test_db?sslmode=disable"
   make test-all
   ```

3. **Setup and Run E2E Tests**:
   ```bash
   make setup-e2e
   make test-e2e
   ```

## Test Types and Coverage

### 1. Unit Tests ⚡ (Fast - ~100ms)
**Purpose**: Test business logic in isolation
**Location**: `*_test.go` files alongside source code
**Command**: `make test-unit`

**What's Tested**:
- Horse age calculation and validation
- Gender enum logic and string representation  
- Date parsing utilities
- HTML path generation
- Image path construction

**Benefits**:
- Extremely fast feedback loop
- No external dependencies
- Easy to debug
- High code coverage of pure functions

### 2. Integration Tests 🔗 (Medium - ~1-5s)
**Purpose**: Test database operations with real database
**Location**: `tests/integration/`
**Command**: `make test-integration`
**Requires**: `TEST_DATABASE_URL` environment variable

**What's Tested**:
- User CRUD operations
- Farm creation and retrieval
- Horse database operations
- Dashboard statistics queries
- Database constraints and relationships

**Benefits**:
- Validates SQL queries work correctly
- Tests database schema and migrations
- Catches database-specific issues
- Verifies data integrity

### 3. HTTP Handler Tests 🌐 (Medium - ~100-500ms)
**Purpose**: Test HTTP endpoints with mocked dependencies
**Location**: `tests/handlers/`
**Command**: `make test-handlers`

**What's Tested**:
- Request parsing (JSON, form data, URL params)
- Input validation and error handling
- HTTP status codes and response format
- Authentication middleware logic
- Route registration and routing

**Benefits**:
- Fast execution with mocked dependencies
- Comprehensive error case testing
- API contract validation
- No external service dependencies

### 4. End-to-End Tests 🎭 (Slow - ~5-30s per test)
**Purpose**: Test complete user workflows in real browser
**Location**: `tests/e2e/`
**Command**: `make test-e2e`
**Requires**: Node.js, test database, Playwright

**What's Tested**:
- Complete authentication flows
- Farm and horse management workflows  
- HTML rendering and template logic
- Navigation and user interactions
- Cross-browser compatibility

**Benefits**:
- Validates actual user experience
- Tests complete system integration
- Catches UI/UX issues
- Verifies JavaScript functionality

## Testing Environment Setup

### Prerequisites
- Go 1.25+
- Node.js 18+ (for E2E tests)
- Docker (for test database)
- Make (for running commands)

### Database Setup Options

#### Option 1: Docker (Recommended)
```bash
# Start test database
make setup-test-db

# Set environment variable
export TEST_DATABASE_URL="postgresql://root@localhost:26258/test_db?sslmode=disable"
```

#### Option 2: Existing CockroachDB
```bash
# Use your existing CockroachDB instance
export TEST_DATABASE_URL="your-test-database-connection-string"

# Run migrations
MIGRATIONS_DSN="$TEST_DATABASE_URL" task db-migrate
```

#### Option 3: CockroachDB Cloud
```bash
# Use CockroachDB Cloud test cluster
export TEST_DATABASE_URL="postgresql://user:pass@host:port/test_db?sslmode=require"

# Run migrations
MIGRATIONS_DSN="$TEST_DATABASE_URL" task db-migrate
```

### E2E Test Setup
```bash
# Install dependencies and browsers
make setup-e2e

# Optional: Set test Stytch credentials
export TEST_STYTCH_PROJECT_ID="your-test-project-id"
export TEST_STYTCH_SECRET="your-test-secret"
```

## Running Tests

### Individual Test Types
```bash
# Unit tests only (fastest)
make test-unit

# Integration tests only
make test-integration

# Handler tests only  
make test-handlers

# E2E tests only
make test-e2e
```

### Combined Test Runs
```bash
# All Go tests (unit + integration + handlers)
make test-all

# Everything including E2E
make test-unit && make test-integration && make test-handlers && make test-e2e
```

### Development Workflows
```bash
# Watch for changes and re-run unit tests
make watch-unit

# Watch for changes and re-run integration tests
make watch-integration

# Run E2E tests with visible browser (for debugging)
make test-e2e-headed

# Open Playwright UI for test development
make test-e2e-ui
```

### Coverage Reports
```bash
# Generate coverage report
make coverage

# View coverage in browser
open coverage.html
```

## Test Data Management

### Unit Tests
- Use hard-coded test data
- No cleanup required
- Focus on edge cases and boundary conditions

### Integration Tests
- Use `TestFixtures` helper for data creation
- Automatic cleanup with registered cleanup functions
- Database transactions for isolation

### Handler Tests
- Use mock database and services
- Predictable responses for different scenarios
- Test both success and error cases

### E2E Tests
- Database seeding before test runs
- Reset to known state between tests
- Use unique identifiers to avoid conflicts

## Authentication Testing

### Unit/Integration Tests
```go
// Use mock authentication
authHelper := testutil.NewAuthTestHelper(db)
sessionToken, user, err := authHelper.CreateAuthenticatedUser(ctx, "Test User", "test@example.com")

// Use in tests
client := testutil.NewTestClient(app).WithAuth(sessionToken)
response := client.Get("/protected/route")
```

### E2E Tests
Currently uses placeholder auth. Options for implementation:

1. **Stytch Test Environment**: Use test project credentials
2. **Authentication Bypass**: Add test-only auth routes
3. **Mock Service**: Intercept Stytch API calls

## Continuous Integration

### GitHub Actions
The project includes a comprehensive CI pipeline:
- **Unit Tests**: Run on every push/PR
- **Integration Tests**: With CockroachDB service
- **Handler Tests**: Fast API testing
- **E2E Tests**: Full browser testing
- **Security Scanning**: Gosec security analysis
- **Linting**: Code quality checks

### Local CI Simulation
```bash
# Run the same tests as CI
make ci-test

# Setup CI environment locally
make ci-setup
```

## Troubleshooting

### Common Issues

#### Database Connection Errors
```bash
# Check if database is running
curl -f http://localhost:26258/health

# Verify connection string
echo $TEST_DATABASE_URL

# Check if migrations ran
psql "$TEST_DATABASE_URL" -c "\dt"
```

#### E2E Test Failures
```bash
# Run with visible browser
make test-e2e-headed

# Check test output
cat tests/e2e/test-results/*/stdout

# View screenshots/videos
ls tests/e2e/test-results/
```

#### Import/Module Errors
```bash
# Clean module cache
go clean -modcache

# Tidy dependencies  
go mod tidy

# Verify imports
go mod verify
```

### Performance Issues
```bash
# Run benchmarks
make bench

# Profile tests
go test -cpuprofile cpu.prof -memprofile mem.prof ./...

# Analyze profiles
go tool pprof cpu.prof
```

## Best Practices

### Writing Tests
1. **Follow AAA Pattern**: Arrange, Act, Assert
2. **Use Table-Driven Tests**: For multiple scenarios
3. **Test Edge Cases**: Empty inputs, boundary values, error conditions
4. **Descriptive Names**: `TestCreateHorse_InvalidGender_ReturnsError`
5. **Independent Tests**: Each test should be isolated

### Test Organization
1. **Co-locate Unit Tests**: Keep `*_test.go` files next to source
2. **Group Integration Tests**: Separate directory with setup/teardown
3. **Shared Test Utilities**: Common helpers in `tests/testutil/`
4. **Clear Documentation**: README in each test directory

### Performance
1. **Run Unit Tests First**: Fast feedback during development
2. **Parallel Integration Tests**: Use `t.Parallel()` when possible
3. **Efficient E2E Tests**: Group related actions, minimize browser restarts
4. **CI Optimization**: Cache dependencies, parallel jobs

### Maintenance
1. **Keep Tests Updated**: When changing code, update tests
2. **Review Test Coverage**: Aim for >80% coverage on critical paths
3. **Remove Obsolete Tests**: Clean up when removing features
4. **Monitor Test Performance**: Watch for slow tests

## Getting Help

- **View Available Commands**: `make help`
- **Test Documentation**: Check README files in test directories
- **Example Usage**: `make example-test-setup`
- **CI Logs**: Check GitHub Actions for detailed error information

This testing strategy provides comprehensive coverage while maintaining fast feedback loops for development. Start with unit tests for immediate feedback, then add integration and E2E tests for complete confidence in your changes.