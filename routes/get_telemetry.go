package routes

import (
	"fmt"

	"github.com/games4l/telemetry-service/providers"
	"github.com/gofiber/fiber/v2"
)

func GetTelemetryUnit(ts *providers.TelemetryService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		item, err := ts.FindById(idParam)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": fmt.Sprintf("telemetry registry %s could not be found", idParam),
			})
		}

		return c.JSON(fiber.Map{
			"message": "registry found",
			"data":    item,
		})
	}
}
