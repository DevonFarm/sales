package testutil

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// mustParseDate parses a date string or panics (for test data creation)
func mustParseDate(dateStr string) time.Time {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		panic(fmt.Sprintf("failed to parse date %s: %v", dateStr, err))
	}
	return date
}

// JSONRequest creates an HTTP request with JSON body
func JSONRequest(method, url string, body interface{}) *httptest.Request {
	var bodyReader io.Reader
	
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			panic(fmt.Sprintf("failed to marshal JSON: %v", err))
		}
		bodyReader = strings.NewReader(string(jsonData))
	}
	
	req := httptest.NewRequest(method, url, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	
	return req
}

// FormRequest creates an HTTP request with form data
func FormRequest(method, url string, formData map[string]string) *httptest.Request {
	values := make([]string, 0, len(formData))
	for key, value := range formData {
		values = append(values, fmt.Sprintf("%s=%s", key, value))
	}
	
	body := strings.Join(values, "&")
	req := httptest.NewRequest(method, url, strings.NewReader(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	
	return req
}

// ParseJSONResponse parses a JSON response body into a struct
func ParseJSONResponse(resp *httptest.Response, dest interface{}) error {
	defer resp.Body.Close()
	return json.NewDecoder(resp.Body).Decode(dest)
}

// GetResponseBody reads the entire response body as a string
func GetResponseBody(resp *httptest.Response) (string, error) {
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// AssertStatusCode checks if response has expected status code
func AssertStatusCode(resp *httptest.Response, expected int) error {
	if resp.StatusCode != expected {
		body, _ := GetResponseBody(resp)
		return fmt.Errorf("expected status %d, got %d. Body: %s", expected, resp.StatusCode, body)
	}
	return nil
}

// AssertContains checks if response body contains expected string
func AssertContains(resp *httptest.Response, expected string) error {
	body, err := GetResponseBody(resp)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}
	
	if !strings.Contains(body, expected) {
		return fmt.Errorf("response body does not contain '%s'. Body: %s", expected, body)
	}
	
	return nil
}

// CreateTestApp creates a Fiber app configured for testing
func CreateTestApp() *fiber.App {
	return fiber.New(fiber.Config{
		// Disable startup message for cleaner test output
		DisableStartupMessage: true,
		// Return errors as JSON in tests
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{
				"error": err.Error(),
			})
		},
	})
}

// WithCookie adds a cookie to the request
func WithCookie(req *httptest.Request, name, value string) {
	req.AddCookie(&http.Cookie{
		Name:  name,
		Value: value,
	})
}

// TestResponse wraps httptest.Response with helper methods
type TestResponse struct {
	*httptest.Response
}

// NewTestResponse wraps a response with helper methods
func NewTestResponse(resp *httptest.Response) *TestResponse {
	return &TestResponse{Response: resp}
}

// AssertStatus checks status code and returns self for chaining
func (r *TestResponse) AssertStatus(expected int) *TestResponse {
	if r.StatusCode != expected {
		panic(fmt.Sprintf("expected status %d, got %d", expected, r.StatusCode))
	}
	return r
}

// AssertJSON parses response as JSON into destination
func (r *TestResponse) AssertJSON(dest interface{}) *TestResponse {
	if err := ParseJSONResponse(r.Response, dest); err != nil {
		panic(fmt.Sprintf("failed to parse JSON response: %v", err))
	}
	return r
}

// AssertBodyContains checks if body contains expected string
func (r *TestResponse) AssertBodyContains(expected string) *TestResponse {
	if err := AssertContains(r.Response, expected); err != nil {
		panic(err.Error())
	}
	return r
}

// TestClient provides a high-level testing client
type TestClient struct {
	app *fiber.App
}

// NewTestClient creates a new test client
func NewTestClient(app *fiber.App) *TestClient {
	return &TestClient{app: app}
}

// Get performs a GET request
func (c *TestClient) Get(url string) *TestResponse {
	req := httptest.NewRequest("GET", url, nil)
	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// Post performs a POST request with JSON body
func (c *TestClient) Post(url string, body interface{}) *TestResponse {
	req := JSONRequest("POST", url, body)
	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// PostForm performs a POST request with form data
func (c *TestClient) PostForm(url string, formData map[string]string) *TestResponse {
	req := FormRequest("POST", url, formData)
	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// Put performs a PUT request with JSON body
func (c *TestClient) Put(url string, body interface{}) *TestResponse {
	req := JSONRequest("PUT", url, body)
	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// Delete performs a DELETE request
func (c *TestClient) Delete(url string) *TestResponse {
	req := httptest.NewRequest("DELETE", url, nil)
	resp, err := c.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// WithAuth adds authentication cookie to subsequent requests
func (c *TestClient) WithAuth(sessionToken string) *AuthenticatedTestClient {
	return &AuthenticatedTestClient{
		client:       c,
		sessionToken: sessionToken,
	}
}

// AuthenticatedTestClient wraps TestClient with authentication
type AuthenticatedTestClient struct {
	client       *TestClient
	sessionToken string
}

// Get performs authenticated GET request
func (c *AuthenticatedTestClient) Get(url string) *TestResponse {
	req := httptest.NewRequest("GET", url, nil)
	WithCookie(req, "stytch_session_token", c.sessionToken)
	resp, err := c.client.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}

// Post performs authenticated POST request
func (c *AuthenticatedTestClient) Post(url string, body interface{}) *TestResponse {
	req := JSONRequest("POST", url, body)
	WithCookie(req, "stytch_session_token", c.sessionToken)
	resp, err := c.client.app.Test(req)
	if err != nil {
		panic(fmt.Sprintf("request failed: %v", err))
	}
	return NewTestResponse(resp)
}