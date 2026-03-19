package util

import (
	"context"
	"fmt"
	"time"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/user"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// MockStytchAuth provides a mock authentication service for testing
type MockStytchAuth struct {
	sessions map[string]*MockSession
	users    map[string]*user.User
}

// MockSession represents a mock Stytch session
type MockSession struct {
	SessionToken string
	UserID       string
	ExpiresAt    time.Time
}

// NewMockStytchAuth creates a new mock authentication service
func NewMockStytchAuth() *MockStytchAuth {
	return &MockStytchAuth{
		sessions: make(map[string]*MockSession),
		users:    make(map[string]*user.User),
	}
}

// CreateMockSession creates a mock session for testing
func (m *MockStytchAuth) CreateMockSession(userID string, user *user.User) string {
	sessionToken := fmt.Sprintf("mock-session-%s-%d", userID, time.Now().Unix())

	m.sessions[sessionToken] = &MockSession{
		SessionToken: sessionToken,
		UserID:       userID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	if user != nil {
		m.users[userID] = user
	}

	return sessionToken
}

// MockAuthenticate simulates Stytch session authentication
func (m *MockStytchAuth) MockAuthenticate(sessionToken string) (*MockAuthResponse, error) {
	session, exists := m.sessions[sessionToken]
	if !exists {
		return nil, fmt.Errorf("invalid session token")
	}

	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	return &MockAuthResponse{
		SessionToken: session.SessionToken,
		Session: &MockSessionData{
			UserID: session.UserID,
		},
	}, nil
}

// MockAuthResponse mimics Stytch authentication response
type MockAuthResponse struct {
	SessionToken string
	Session      *MockSessionData
}

// MockSessionData mimics Stytch session data
type MockSessionData struct {
	UserID string
}

// RequireAuth creates a mock authentication middleware for testing
func (m *MockStytchAuth) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("stytch_session_token")
		if token == "" {
			return c.Redirect("/login")
		}

		authResp, err := m.MockAuthenticate(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "unauthorized"})
		}

		// Set user info in locals
		c.Locals("stytch_session", authResp.Session)
		c.Locals("stytch_user_id", authResp.Session.UserID)

		return c.Next()
	}
}

// AuthTestHelper provides utilities for authentication testing
type AuthTestHelper struct {
	mockAuth *MockStytchAuth
	db       *database.DB
}

// NewAuthTestHelper creates a new authentication test helper
func NewAuthTestHelper(db *database.DB) *AuthTestHelper {
	return &AuthTestHelper{
		mockAuth: NewMockStytchAuth(),
		db:       db,
	}
}

// CreateAuthenticatedUser creates a user and returns session token for testing
func (h *AuthTestHelper) CreateAuthenticatedUser(ctx context.Context, name, email string) (string, *user.User, error) {
	stytchID := fmt.Sprintf("stytch-%s", uuid.New().String())

	// Create user in database
	testUser, err := user.NewUser(ctx, h.db, name, email, stytchID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to create test user: %w", err)
	}

	// Create mock session
	sessionToken := h.mockAuth.CreateMockSession(stytchID, testUser)

	return sessionToken, testUser, nil
}

// GetMockAuth returns the mock authentication service
func (h *AuthTestHelper) GetMockAuth() *MockStytchAuth {
	return h.mockAuth
}

// SetAuthenticatedContext sets up Fiber context with authenticated user
func (h *AuthTestHelper) SetAuthenticatedContext(c *fiber.Ctx, userID string) {
	c.Locals("stytch_user_id", userID)
	c.Locals("stytch_session", &MockSessionData{UserID: userID})
}

// CreateTestAuthMiddleware creates middleware that bypasses auth for testing
func CreateTestAuthMiddleware(userID string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set test user context
		c.Locals("stytch_user_id", userID)
		c.Locals("stytch_session", &MockSessionData{UserID: userID})
		return c.Next()
	}
}

// WithAuthentication runs a test function with an authenticated context
func WithAuthentication(userID string, fn func(c *fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set authentication context
		c.Locals("stytch_user_id", userID)
		c.Locals("stytch_session", &MockSessionData{UserID: userID})

		// Run the test function
		return fn(c)
	}
}

// MockStytchClient provides a mock Stytch client for testing
type MockStytchClient struct {
	magicLinks map[string]*MockMagicLink
	sessions   map[string]*MockSession
}

// MockMagicLink represents a mock magic link
type MockMagicLink struct {
	Token  string
	Email  string
	UserID string
}

// NewMockStytchClient creates a new mock Stytch client
func NewMockStytchClient() *MockStytchClient {
	return &MockStytchClient{
		magicLinks: make(map[string]*MockMagicLink),
		sessions:   make(map[string]*MockSession),
	}
}

// MockSendMagicLink simulates sending a magic link
func (m *MockStytchClient) MockSendMagicLink(email string) (*MockMagicLink, error) {
	userID := fmt.Sprintf("user-%s", uuid.New().String())
	token := fmt.Sprintf("magic-token-%d", time.Now().Unix())

	link := &MockMagicLink{
		Token:  token,
		Email:  email,
		UserID: userID,
	}

	m.magicLinks[token] = link

	return link, nil
}

// MockAuthenticateMagicLink simulates magic link authentication
func (m *MockStytchClient) MockAuthenticateMagicLink(token string) (*MockAuthResponse, error) {
	link, exists := m.magicLinks[token]
	if !exists {
		return nil, fmt.Errorf("invalid magic link token")
	}

	sessionToken := fmt.Sprintf("session-%s-%d", link.UserID, time.Now().Unix())

	session := &MockSession{
		SessionToken: sessionToken,
		UserID:       link.UserID,
		ExpiresAt:    time.Now().Add(24 * time.Hour),
	}

	m.sessions[sessionToken] = session

	return &MockAuthResponse{
		SessionToken: sessionToken,
		Session: &MockSessionData{
			UserID: link.UserID,
		},
	}, nil
}
