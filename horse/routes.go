package horse

import (
	"github.com/DevonFarm/sales/auth"
	"github.com/DevonFarm/sales/database"
	"github.com/DevonFarm/sales/farm"
	"github.com/DevonFarm/sales/utils"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(app *fiber.App, db *database.DB, auth *auth.StytchAuth) {
	farmGroup := app.Group("/farm/:farmID", auth.RequireAuth())
	farmGroup.Get("/", getDashboard(db))
	farmGroup.Get("/horses", getHorses(db))
	farmGroup.Get("/horse/:id", getHorse(db))
	farmGroup.Post("/horse", createHorse(db))
	farmGroup.Put("/horse/:id", updateHorse(db))
	farmGroup.Delete("/horse/:id", deleteHorse(db))

	app.Get("/list", func(c *fiber.Ctx) error {
		// TODO: need a different template to list horses
		return c.Render("templates/create", fiber.Map{
			"Title": "Listing",
		})
	})
}

func getDashboard(db *database.DB) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		farmID := c.Params("farmID")
		if farmID == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "farm ID is required"})
		}

		// Get farm details
		f, err := farm.GetFarm(c.Context(), db, farmID)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "farm not found"})
		}

		horses, err := GetHorsesByFarmID(c.Context(), db, f.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get horses"})
		}

		// Get dashboard statistics
		stats, err := GetDashboardStats(c.Context(), db, f.ID)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get stats"})
		}

		return c.Render("templates/dashboard", fiber.Map{
			"Title":  f.Name + " Dashboard",
			"Farm":   f,
			"Horses": horses,
			"Stats":  stats,
		})
	}
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
		if err := h.Save(c.Context(), db); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.Status(fiber.StatusCreated).JSON(h)
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
