package routes

import (
	"strings"

	"github.com/games4l/telemetry-service/providers"
	"github.com/gofiber/fiber/v2"
)

func GetTelemetryUnit(ts *providers.TelemetryService, ap *providers.AuthProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		item, err := ts.FindById(idParam)

		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "telemetry registry %s could not be found " + idParam,
			})
		}

		authHeaderS := strings.Split(c.Get("Authorization"), " ")

		if len(authHeaderS) < 3 {
			item.PacientName = "<OMITTED>"

			return c.JSON(fiber.Map{
				"message": "registry found",
				"data":    item,
			})
		}

		if authHeaderS[0] != "Signature" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid auth strategy " + authHeaderS[0],
			})
		}

		encodingS := providers.ByteEncoding(authHeaderS[1])

		if encodingS != providers.ByteEncodingBase64 && encodingS != providers.ByteEncodingHex {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid encoding strategy " + authHeaderS[1],
			})
		}

		vErr := ap.ValidateSignature(encodingS, c.Body(), []byte(authHeaderS[2]))

		if vErr != nil {
			return c.Status(vErr.Status()).JSON(fiber.Map{
				"error": vErr.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "registry found",
			"data":    item,
		})
	}
}
