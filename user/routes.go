package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/DevonFarm/sales/database"
)

func RegisterRoutes(app *fiber.App, db *database.DB, authMiddleware fiber.Handler) {
	userGroup := app.Group("/user/:id", authMiddleware)
	userGroup.Get("/profile", getProfile(db))
	userGroup.Post("/profile", updateProfile(db))
}

func getProfile(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user ID is required"})
		}

		user, err := GetUser(c.Context(), db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user"})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		return c.Render("templates/profile", fiber.Map{
			"Title": "Edit Profile",
			"User":  user,
		})
	}
}

func updateProfile(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		userID := c.Params("id")
		if userID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "user ID is required"})
		}

		// Get existing user
		user, err := GetUser(c.Context(), db, userID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get user"})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "user not found"})
		}

		// Parse form data
		if err := c.BodyParser(user); err != nil {
			return c.Status(fiber.StatusBadRequest).Render("templates/profile", fiber.Map{
				"Title": "Edit Profile",
				"User":  user,
				"Error": "Invalid form data",
			})
		}

		// Save updated user
		if err := user.Save(c.Context(), db); err != nil {
			return c.Status(fiber.StatusInternalServerError).Render("templates/profile", fiber.Map{
				"Title": "Edit Profile",
				"User":  user,
				"Error": "Failed to update profile",
			})
		}

		// Redirect to farm dashboard if user has a farm, otherwise to farm creation
		if user.FarmID == uuid.Nil {
			return c.Redirect(fmt.Sprintf("/new/farm/%s", user.ID))
		}
		return c.Redirect(fmt.Sprintf("/farm/%s", user.FarmID))
	}
}
