package utils

import (
	"github.com/DevonFarm/sales/database"

	"github.com/gofiber/fiber/v2"
)

type FiberHandlerWithDB func(db *database.DB) func(*fiber.Ctx) error
