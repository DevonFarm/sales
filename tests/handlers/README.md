# HTTP Handler Tests

These tests verify HTTP endpoints and request/response handling with mocked dependencies.

## What These Tests Cover

- **Request Parsing**: JSON body parsing, form data, URL parameters
- **Validation**: Input validation and error handling
- **Response Format**: Correct HTTP status codes and response bodies
- **Error Cases**: Invalid input, malformed requests, constraint violations

## Testing Strategy

### Mocked Dependencies
- **Database**: Custom mock that implements expected database interface
- **Authentication**: Middleware can be bypassed or mocked
- **External Services**: Stytch and other external calls are mocked

### Test Structure
```go
func TestHandlerName_Scenario(t *testing.T) {
    // 1. Setup mock dependencies
    mockDB := NewMockDB()
    
    // 2. Create Fiber app with test routes
    app := fiber.New()
    app.Post("/route", handlerFunction)
    
    // 3. Create test request
    req := httptest.NewRequest("POST", "/route", requestBody)
    
    // 4. Execute request
    resp, err := app.Test(req)
    
    // 5. Assert response
    assert.Equal(t, expectedStatus, resp.StatusCode)
    // ... more assertions
}
```

## Running Tests

```bash
go test ./tests/handlers/ -v
```

## Benefits

1. **Fast**: No database or external service calls
2. **Isolated**: Each test is independent
3. **Predictable**: Mocked responses are consistent
4. **Comprehensive**: Can test error conditions easily

## Mock Database

The `MockDB` struct provides:
- In-memory storage for test data
- Implements expected database interface methods
- Allows testing of database error conditions
- Provides predictable responses for different scenarios

## Authentication Testing

For routes that require authentication:
1. Mock the authentication middleware
2. Set expected user context in Fiber locals
3. Test both authenticated and unauthenticated scenarios

Example:
```go
// Mock authentication middleware
app.Use("/protected/*", func(c *fiber.Ctx) error {
    // Set mock user data
    c.Locals("stytch_user_id", "test-user-123")
    return c.Next()
})
```