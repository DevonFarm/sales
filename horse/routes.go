package horse

import (
	"github.com/DevonFarm/sales/database"

	"github.com/gofiber/fiber/v2"
)

func Routes(app *fiber.App, db *database.DB) {
	app.Get("/horses", getHorses(db))
	app.Get("/horses/:id", getHorse(db))
	app.Post("/horses", createHorse(db))
	app.Put("/horses/:id", updateHorse(db))
	app.Delete("/horses/:id", deleteHorse(db))

	app.Get("/list", func(c *fiber.Ctx) error {
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
