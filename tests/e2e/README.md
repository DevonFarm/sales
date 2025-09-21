# End-to-End Tests with Playwright

These tests verify complete user workflows through a real browser.

## Setup

1. **Install Dependencies**:
   ```bash
   cd tests/e2e
   npm install
   npx playwright install
   ```

2. **Environment Setup**:
   Create `.env` file in the project root with test environment variables:
   ```bash
   TEST_DATABASE_URL="postgresql://root@localhost:26257/test_db?sslmode=disable"
   TEST_STYTCH_PROJECT_ID="your-test-project-id"
   TEST_STYTCH_SECRET="your-test-secret"
   ```

3. **Database Setup**:
   Ensure your test database is running and migrations are applied:
   ```bash
   MIGRATIONS_DSN="$TEST_DATABASE_URL" task db-migrate
   ```

## Running Tests

```bash
# Run all tests
npm test

# Run tests with browser UI visible
npm run test:headed

# Debug tests interactively
npm run test:debug

# Open Playwright UI for test development
npm run test:ui
```

## Test Structure

### Homepage Tests (`homepage.spec.ts`)
- Basic page loading and navigation
- Static content verification
- Link functionality

### Authentication Tests (`auth.spec.ts`)
- Login form display and validation
- Magic link request flow
- Logout functionality
- **Note**: Full auth testing requires Stytch test environment setup

### Farm Management Tests (`farm-management.spec.ts`)
- Farm creation workflow
- Dashboard display and statistics
- Horse CRUD operations through UI
- Form validation
- **Note**: Currently skipped pending auth setup

## Authentication Strategy

The tests currently use placeholder authentication. To implement full auth testing:

### Option 1: Stytch Test Environment
```typescript
// Use Stytch test project credentials
const testAuth = {
  projectId: 'test-project-live-12345',
  secret: 'secret-test-67890'
};

// Send magic link to test email
await page.fill('input[name="email"]', 'test@stytch.com');
// Intercept email and extract magic link
// Navigate to callback URL
```

### Option 2: Authentication Bypass
```go
// Add test route in Go app when TEST_MODE=true
if os.Getenv("TEST_MODE") == "true" {
    app.Get("/test/auth/:userID", func(c *fiber.Ctx) error {
        // Set test session cookie
        c.Cookie(&fiber.Cookie{
            Name:  "stytch_session_token", 
            Value: "test-token-" + c.Params("userID"),
        })
        return c.Redirect("/")
    })
}
```

### Option 3: Mock Authentication Service
```typescript
// Intercept Stytch API calls
await page.route('**/stytch.com/**', route => {
  route.fulfill({
    status: 200,
    body: JSON.stringify({ success: true, token: 'mock-token' })
  });
});
```

## Test Data Management

### Database Seeding
```typescript
// Before each test
test.beforeEach(async ({ page }) => {
  // Reset database to known state
  await page.request.post('/test/reset-db');
  
  // Seed with test data
  await page.request.post('/test/seed', {
    data: testData
  });
});
```

### Test Isolation
- Each test should clean up its data
- Use unique test data identifiers
- Consider database transactions that rollback

## CI/CD Integration

Example GitHub Actions workflow:
```yaml
- name: Run E2E Tests
  run: |
    # Start test database
    docker run -d --name test-db cockroachdb/cockroach:latest
    
    # Set environment variables
    export TEST_DATABASE_URL="..."
    
    # Run migrations
    task db-migrate
    
    # Run tests
    cd tests/e2e
    npm ci
    npx playwright install --with-deps
    npm test
```

## Benefits of E2E Testing

1. **User Experience Validation**: Tests actual user workflows
2. **Integration Testing**: Verifies all components work together
3. **Visual Regression**: Can detect UI changes
4. **Cross-Browser Testing**: Ensures compatibility
5. **Real Environment**: Tests with actual database and services

## Current Limitations

1. **Authentication**: Requires Stytch test environment setup
2. **Test Data**: Needs database seeding strategy
3. **Email Testing**: Magic link flow needs email interception
4. **Performance**: Slower than unit/integration tests

## Next Steps

1. Set up Stytch test environment
2. Implement authentication bypass for tests
3. Add database seeding utilities
4. Create page object models for maintainability
5. Add visual regression testing