# Integration Tests

These tests verify database operations and data integrity.

## Setup

1. **Test Database**: You need a separate test database instance
   ```bash
   export TEST_DATABASE_URL="postgresql://username:password@host:port/test_database"
   ```

2. **Run Migrations**: Apply migrations to your test database first
   ```bash
   # Using your existing migration setup
   MIGRATIONS_DSN="$TEST_DATABASE_URL" task db-migrate
   ```

3. **Run Tests**:
   ```bash
   go test ./tests/integration/ -v
   ```

## What These Tests Cover

- **User CRUD Operations**: Creating, reading, updating users
- **Farm CRUD Operations**: Farm creation and retrieval  
- **Horse CRUD Operations**: Horse creation, retrieval, and farm associations
- **Dashboard Statistics**: Aggregate queries for horse counts by gender
- **Database Constraints**: Testing validation and constraint enforcement

## Test Database Strategy

- Tests use a separate `TEST_DATABASE_URL` environment variable
- Each test cleans up its data using `cleanupTestData()`
- Tests are designed to be independent and can run in any order
- If no test database is configured, tests are skipped

## CI/CD Integration

For CI/CD pipelines, you can:
1. Spin up a temporary CockroachDB instance
2. Run migrations
3. Set `TEST_DATABASE_URL` 
4. Run integration tests
5. Tear down the database

Example with Docker:
```bash
# Start test database
docker run -d --name test-cockroach -p 26257:26257 cockroachdb/cockroach:latest start-single-node --insecure

# Wait for startup, run migrations, then test
export TEST_DATABASE_URL="postgresql://root@localhost:26257/defaultdb?sslmode=disable"
go test ./tests/integration/ -v

# Cleanup
docker stop test-cockroach && docker rm test-cockroach
```