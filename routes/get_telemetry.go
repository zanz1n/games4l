package routes

import (
	"fmt"
	"strings"

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

		authHeaderS := strings.Split(c.Get("Authorization"), " ")

		if len(authHeaderS) < 2 {
			item.PacientName = "<OMITTED>"

			return c.JSON(fiber.Map{
				"message": "registry found",
				"data":    item,
			})
		}

		if authHeaderS[0] != "Bearer" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": fmt.Sprintf("invalid authorization stategy %s", authHeaderS[0]),
			})
		}

		return c.JSON(fiber.Map{
			"message": "registry found",
			"data":    item,
		})
	}
}
