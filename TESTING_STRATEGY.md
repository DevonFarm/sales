# Testing Strategy for Devon Farm Sales

## Overview
This Go web application presents unique testing challenges due to its heavy integration with:
- CockroachDB database operations
- Stytch authentication service  
- HTML template rendering
- Complex user workflows

## Testing Pyramid

### 1. Unit Tests (Go) - Fast & Isolated
**Purpose**: Test pure business logic without external dependencies
**Location**: `*_test.go` files alongside source code
**Coverage**:
- Model methods (`Horse.Age()`, validation logic)
- Utility functions (`utils/date.go`, `utils/error.go`) 
- Enum validation (`horse/gender.go`)
- Template path generation
- Input parsing and validation

### 2. Integration Tests (Go) - Database Operations
**Purpose**: Test database interactions with real database
**Location**: `tests/integration/`
**Coverage**:
- CRUD operations for all models
- Database migrations
- Query correctness and edge cases
- Transaction handling
- Database constraints and relationships

### 3. HTTP Handler Tests (Go) - API Layer
**Purpose**: Test HTTP handlers with mocked dependencies
**Location**: `tests/handlers/`
**Coverage**:
- Route registration and routing
- Request/response handling
- Authentication middleware
- Form parsing and validation
- Error handling and status codes
- Template rendering (without full HTML validation)

### 4. End-to-End Tests (Playwright) - Full User Flows
**Purpose**: Test complete user workflows with real browser
**Location**: `tests/e2e/`
**Coverage**:
- Complete authentication flows (magic link, login, logout)
- Farm creation and management
- Horse CRUD operations through UI
- Navigation and redirects
- HTML rendering and template logic
- Session management and cookies

## Test Database Strategy

### For Integration Tests
- Use separate test CockroachDB database
- Run migrations before each test suite
- Clean up data between tests
- Use transactions that rollback for test isolation

### For E2E Tests  
- Use dedicated test database instance
- Seed with known test data
- Reset to known state between test runs
- May share database across E2E tests for performance

## Authentication Testing Strategy

### Unit/Integration Tests
- Mock Stytch API responses
- Test authentication middleware logic
- Validate session handling

### E2E Tests
- Use Stytch test environment
- Create test user accounts
- Test complete magic link flows
- Validate session persistence

## Implementation Priority

1. **Start with Unit Tests** - Quick wins, build testing habits
2. **Add Integration Tests** - Validate database layer reliability  
3. **HTTP Handler Tests** - Ensure API contracts work
4. **Playwright E2E Tests** - Validate complete user experience

## Tools and Dependencies

### Go Testing
- Standard `testing` package
- `testify` for assertions and mocking
- `pgx` test utilities for database testing
- `httptest` for HTTP handler testing

### Playwright
- TypeScript/JavaScript Playwright setup
- Page Object Model for maintainable tests
- Test data management utilities
- Screenshot/video capture on failures

## Benefits of This Approach

1. **Fast Feedback Loop**: Unit tests run in milliseconds
2. **Reliable Database Testing**: Integration tests catch SQL issues
3. **API Contract Validation**: Handler tests ensure HTTP layer works
4. **User Experience Validation**: E2E tests validate the actual user flows
5. **Maintainable**: Clear separation of concerns and test types
6. **CI/CD Friendly**: Different test types can run at different stages

## Getting Started

1. Set up basic unit tests for utility functions
2. Configure test database and integration test framework
3. Add HTTP handler tests with mocked dependencies
4. Set up Playwright for critical user flows
5. Gradually increase coverage in each layer