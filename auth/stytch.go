package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks/email"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/stytchapi"
)

const defaultCookieName = "stytch_session_jwt"

type StytchAuth struct {
	CookieName string
	Client     *stytchapi.API
}

// NewStytchFromEnv creates a Stytch client from environment variables:
// STYTCH_PROJECT_ID, STYTCH_SECRET
func NewStytchFromEnv() (*StytchAuth, error) {
	projectID := os.Getenv("STYTCH_PROJECT_ID")
	secret := os.Getenv("STYTCH_SECRET")
	if projectID == "" || secret == "" {
		return nil, errors.New("missing STYTCH_PROJECT_ID or STYTCH_SECRET env var")
	}

	client, err := stytchapi.NewClient(projectID, secret)
	if err != nil {
		return nil, fmt.Errorf("stytchapi.NewClient: %w", err)
	}

	return &StytchAuth{
		Client:     client,
		CookieName: defaultCookieName,
	}, nil
}

// Register mounts auth routes: GET /login, POST /login, GET /auth/callback, POST /logout
func (a *StytchAuth) Register(app *fiber.App) {
	app.Get("/login", a.renderLogin)
	app.Post("/login", a.sendMagicLink)
	app.Get("/auth/callback", a.magicLinkCallback)
	app.Post("/logout", a.logout)
}

// RequireAuth verifies the session JWT cookie and sets user info in Locals.
func (a *StytchAuth) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		jwt := c.Cookies(a.CookieName)
		if jwt == "" {
			return c.Redirect("/login")
		}
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()
		// Authenticate the session JWT with Stytch
		res, err := a.Client.Sessions.AuthenticateJWT(ctx, 5*time.Minute, &sessions.AuthenticateParams{SessionJWT: jwt})
		if err != nil {
			c.SendStatus(fiber.StatusUnauthorized)
			return err
		}
		c.Cookie(&fiber.Cookie{Name: a.CookieName, Value: "", Expires: time.Unix(0, 0), HTTPOnly: true, Secure: isSecure(c), SameSite: fiber.CookieSameSiteLaxMode})

		// Stash user info for handlers/templates
		c.Locals("stytch_session", res.Session)
		c.Locals("stytch_user_id", res.Session.UserID)
		return c.Next()
	}
}

func (a *StytchAuth) renderLogin(c *fiber.Ctx) error {
	// If already logged in, skip
	if c.Cookies(a.CookieName) != "" {
		return c.Redirect("/")
	}
	return c.Render("login", fiber.Map{
		"Title": "Log in",
	})
}

func (a *StytchAuth) sendMagicLink(c *fiber.Ctx) error {
	type form struct {
		Email string `form:"email"`
	}
	var f form
	if err := c.BodyParser(&f); err != nil || f.Email == "" {
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"Title": "Log in",
			"Error": "Enter a valid email",
		})
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	// Send the magic link via email
	params := email.LoginOrCreateParams{Email: f.Email}
	if _, err := a.Client.MagicLinks.Email.LoginOrCreate(ctx, &params); err != nil {
		return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
			"Title": "Log in",
			"Error": fmt.Sprintf("Couldn't send link: %v", err),
		})
	}

	return c.Render("login_sent", fiber.Map{"Title": "Check your email", "Email": f.Email})
}

func (a *StytchAuth) magicLinkCallback(c *fiber.Ctx) error {
	token := c.Query("token")
	if token == "" {
		return c.Status(fiber.StatusBadRequest).SendString("missing token")
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()

	res, err := a.Client.MagicLinks.Authenticate(ctx, &magiclinks.AuthenticateParams{
		Token:                  token,
		SessionDurationMinutes: 60,
	})
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).SendString("invalid or expired link")
	}

	// Set the session JWT cookie
	c.Cookie(&fiber.Cookie{
		Name:     a.CookieName,
		Value:    res.SessionJWT,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   isSecure(c),
		SameSite: fiber.CookieSameSiteLaxMode,
	})

	// Redirect to home or to a 'next' param if present
	next := c.Query("next", "/")
	return c.Redirect(next)
}

func (a *StytchAuth) logout(c *fiber.Ctx) error {
	jwt := c.Cookies(a.CookieName)
	if jwt != "" {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		_, _ = a.Client.Sessions.Revoke(ctx, &sessions.RevokeParams{SessionJWT: jwt})
		c.Cookie(&fiber.Cookie{Name: a.CookieName, Value: "", Expires: time.Unix(0, 0), HTTPOnly: true, Secure: isSecure(c), SameSite: fiber.CookieSameSiteLaxMode, Path: "/"})
		return c.Redirect("/")
	}
	return c.Redirect("/")
}

func isSecure(c *fiber.Ctx) bool {
	// Treat X-Forwarded-Proto as signal when behind a proxy
	if strings.EqualFold(string(c.Protocol()), "https") {
		return true
	}
	if p := c.Get("X-Forwarded-Proto"); strings.EqualFold(p, "https") {
		return true
	}
	return false
}
