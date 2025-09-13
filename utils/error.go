package utils

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func LogAndRespondError(c *fiber.Ctx, message string, err error, statusCode int) error {
	fmt.Printf("%s: %v\n", message, err)
	return c.Status(statusCode).SendString(message)
}
