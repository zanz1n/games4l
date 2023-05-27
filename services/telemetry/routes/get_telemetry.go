package routes

import (
	"context"
	"strings"
	"time"

	"github.com/games4l/backend/libs/auth"
	"github.com/games4l/backend/libs/telemetry"
	"github.com/gofiber/fiber/v2"
)

func GetBySimilarName(ts *telemetry.TelemetryService, ap *auth.AuthProvider) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		authHeaderS := strings.Split(c.Get("Authorization"), " ")

		if len(authHeaderS) < 3 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "this route requires admin authorization",
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

		nameParam := c.Query("name")

		if nameParam == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "name query param must be provided",
			})
		}

		ctx, cancel := context.WithTimeout(context.Background(), 32*time.Second)
		defer cancel()

		result, err := ts.FindSimilarName(ctx, nameParam)

		if err != nil {
			return c.Status(err.Status()).JSON(fiber.Map{
				"error": err.Error(),
			})
		}

		return c.JSON(fiber.Map{
			"message": "success",
			"data":    result,
		})
	}
}

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
