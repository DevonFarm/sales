package auth

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/magiclinks/email"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/sessions"
	"github.com/stytchauth/stytch-go/v16/stytch/consumer/stytchapi"

	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/user"
	"github.com/DevonFarm/sales/utils"
)

const defaultCookieName = "stytch_session_token"

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
func (a *StytchAuth) Register(app *fiber.App, db *database.DB) {
	app.Get("/login", a.renderLogin(db))
	app.Post("/login", a.sendMagicLink(db))
	app.Get("/auth/callback", a.magicLinkCallback(db))
	app.Post("/logout", a.logout)
}

// RequireAuth verifies the session token cookie and sets user info in Locals.
func (a *StytchAuth) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies(a.CookieName)
		if token == "" {
			return c.Redirect("/login")
		}
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()
		// Authenticate the session token with Stytch
		res, err := a.Client.Sessions.Authenticate(ctx, &sessions.AuthenticateParams{SessionToken: token})
		if err != nil {
			c.SendStatus(fiber.StatusUnauthorized)
			return err
		}
		// Refresh the cookie with a new expiration time
		c.Cookie(&fiber.Cookie{
			Name:     a.CookieName,
			Value:    res.SessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			Secure:   isSecure(c),
			SameSite: fiber.CookieSameSiteLaxMode,
			Path:     "/",
		})

		// Stash user info for handlers/templates
		c.Locals("stytch_session", res.Session)
		c.Locals("stytch_user_id", res.Session.UserID)
		return c.Next()
	}
}

func (a *StytchAuth) renderLogin(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// If already logged in, skip
		token := c.Cookies(a.CookieName)
		if token != "" {
			// check if session is valid and get farm ID from the user
			res, err := a.Client.Sessions.Authenticate(c.Context(), &sessions.AuthenticateParams{SessionToken: token})
			if err == nil {
				stytchUserID := res.Session.UserID
				u, err := user.GetUserByStytchID(c.Context(), db, stytchUserID)
				if err == nil && u != nil {
					if u.FarmID == uuid.Nil {
						// No farm yet, go to create farm page
						return c.Redirect(fmt.Sprintf("/new/farm/%s", u.ID))
					}
					// Redirect to the user's farm dashboard
					return c.Redirect(fmt.Sprintf("/farm/%s", u.FarmID))
				} else {
					log.Warnf("user not found for stytch ID %s: %v", stytchUserID, err)
				}
			} else {
				log.Warnf("invalid session token: %v", err)
			}
		}
		return c.Render("templates/login", fiber.Map{
			"Title": "Log in",
		})
	}
}

func (a *StytchAuth) sendMagicLink(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var u user.User
		if err := c.BodyParser(&u); err != nil {
			return c.Status(fiber.StatusBadRequest).Render("login", fiber.Map{
				"Title": "Log in",
				"Error": "Enter a valid name and email",
			})
		}

		// Send the magic link via email
		params := email.LoginOrCreateParams{Email: u.Email}
		res, err := a.Client.MagicLinks.Email.LoginOrCreate(c.Context(), &params)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("failed to send magic link")
		}

		_, err = user.NewUser(c.Context(), db, u.Name, u.Email, res.UserID)
		if err != nil {
			return utils.LogAndRespondError(
				c,
				"failed to create user",
				err,
				fiber.StatusInternalServerError,
			)
		}

		return c.Render(
			"templates/login_sent",
			fiber.Map{"Title": "Check your email", "Email": u.Email},
		)
	}
}

func (a *StytchAuth) magicLinkCallback(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
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

		// Set the session token cookie
		c.Cookie(&fiber.Cookie{
			Name:     a.CookieName,
			Value:    res.SessionToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HTTPOnly: true,
			Secure:   isSecure(c),
			SameSite: fiber.CookieSameSiteLaxMode,
			Path:     "/",
		})

		u, err := user.GetUserByStytchID(c.Context(), db, res.UserID)
		if err != nil {
			return utils.LogAndRespondError(
				c,
				"failed to get user by stytch ID",
				err,
				fiber.StatusInternalServerError,
			)
		}
		if u == nil {
			return c.Status(fiber.StatusInternalServerError).SendString("user not found")
		}

		if u.FarmID == uuid.Nil {
			// No farm yet, go to create farm page
			return c.Redirect(fmt.Sprintf("/new/farm/%s", u.ID))
		}
		// Redirect to the user's farm dashboard
		return c.Redirect(fmt.Sprintf("/farm/%s", u.FarmID))
	}
}

func (a *StytchAuth) logout(c *fiber.Ctx) error {
	token := c.Cookies(a.CookieName)
	if token != "" {
		ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
		defer cancel()

		_, _ = a.Client.Sessions.Revoke(ctx, &sessions.RevokeParams{SessionToken: token})
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
