package routes

import (
	"github.com/games4l/backend/libs/telemetry"
	"github.com/gofiber/fiber/v2"
)

func PostTelemetry(ts *telemetry.TelemetryService) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		telemetryData := telemetry.CreateTelemetryUnitData{}

		err := c.BodyParser(&telemetryData)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		result, fErr := ts.Create(&telemetryData)

		if fErr != nil {
			return c.Status(fErr.Status()).JSON(fiber.Map{
				"error": fErr.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"data":    result,
		})
	}
}
