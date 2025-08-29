package horse

import (
	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/server"
	"github.com/DevonFarm/sales/utils"

	"github.com/gofiber/fiber/v2"
)

func Routes(cfg *server.Server) {
	farm := cfg.App.Group("/:farmID", cfg.Auth.RequireAuth())
	farm.Get("/horses", getHorses(cfg.DB))
	farm.Get("/horse/:id", getHorse(cfg.DB))
	farm.Post("/horse", createHorse(cfg.DB))
	farm.Put("/horse/:id", updateHorse(cfg.DB))
	farm.Delete("/horse/:id", deleteHorse(cfg.DB))

	cfg.App.Get("/list", func(c *fiber.Ctx) error {
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
