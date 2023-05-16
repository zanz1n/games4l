package routes

import (
	"github.com/games4l/telemetry-service/providers"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func PostTelemetry(ts *providers.TelemetryService) func(c *fiber.Ctx) error {
	validate := validator.New()

	return func(c *fiber.Ctx) error {
		telemetryData := providers.CreateTelemetryUnitData{}

		c.BodyParser(&telemetryData)

		err := validate.Struct(telemetryData)

		if err != nil {
			return c.JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		result, err := ts.Create(&telemetryData)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "something went wrong",
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"data":    result,
		})
	}
}
