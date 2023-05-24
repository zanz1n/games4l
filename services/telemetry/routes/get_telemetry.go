package routes

import (
	"strings"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/telemetry"
	"github.com/gofiber/fiber/v2"
)

func GetTelemetryUnit(ts *telemetry.TelemetryService, ap *auth.AuthProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id")

		item, fErr := ts.FindById(idParam)

		if fErr != nil {
			return c.Status(fErr.Status()).JSON(fiber.Map{
				"error": fErr.Error(),
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

		encodingS := auth.ByteEncoding(authHeaderS[1])

		if encodingS != auth.ByteEncodingBase64 && encodingS != auth.ByteEncodingHex {
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