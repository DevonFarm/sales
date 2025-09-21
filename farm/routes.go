package farm

import (
	"fmt"

	"github.com/gofiber/fiber/v2"

	"github.com/DevonFarm/sales/auth"
	"github.com/DevonFarm/sales/database"
)

func RegisterRoutes(app *fiber.App, db *database.DB, auth *auth.StytchAuth) {
	newFarm := app.Group("/new/farm", auth.RequireAuth())
	newFarm.Get("/:userID", newFarmForm)
	newFarm.Post("/:userID", createFarm(db))
}

func newFarmForm(c *fiber.Ctx) error {
	return c.Render("templates/new_farm", fiber.Map{
		"Title":  "Create New Farm",
		"UserID": c.Params("userID"),
	})
}

func createFarm(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var f Farm
		if err := c.BodyParser(&f); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		userID := c.Params("userID")
		if err := f.Save(c.Context(), db, userID); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).Redirect(fmt.Sprintf("/farm/%s", f.ID))
	}
}
