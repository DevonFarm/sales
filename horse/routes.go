package horse

import (
	"github.com/DevonFarm/sales/auth"
	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *database.DB, auth *auth.StytchAuth) {
	farm := app.Group("/:farmID", auth.RequireAuth())
	farm.Get("/horses", getHorses(db))
	farm.Get("/horse/:id", getHorse(db))
	farm.Post("/horse", createHorse(db))
	farm.Put("/horse/:id", updateHorse(db))
	farm.Delete("/horse/:id", deleteHorse(db))

	app.Get("/list", func(c *fiber.Ctx) error {
		// TODO: need a different template to list horses
		return c.Render("templates/create", fiber.Map{
			"Title": "Listing",
		})
	})
}

func getHorses(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func getHorse(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func createHorse(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		var h Horse
		if err := c.BodyParser(&h); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}
		dateStr := c.FormValue("date_of_birth")
		if dateStr != "" {
			dob, err := utils.ParseDate(dateStr)
			if err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
			}
			h.DateOfBirth = dob
		}
		h.Save(c.Context(), db)
		return c.SendStatus(fiber.StatusCreated)
	}
}

func updateHorse(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}

func deleteHorse(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		// TODO: implement
		return c.SendStatus(fiber.StatusNotImplemented)
	}
}
