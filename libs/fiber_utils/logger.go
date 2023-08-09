package fiberutils

import (
	"time"

	"github.com/games4l/backend/libs/logger"
	"github.com/gofiber/fiber/v2"
)

var log = logger.NewLogger("http_log")

func statusColor(code int, colors fiber.Colors) string {
	switch {
	case code >= fiber.StatusOK && code < fiber.StatusMultipleChoices:
		return colors.Green
	case code >= fiber.StatusMultipleChoices && code < fiber.StatusBadRequest:
		return colors.Blue
	case code >= fiber.StatusBadRequest && code < fiber.StatusInternalServerError:
		return colors.Yellow
	default:
		return colors.Red
	}
}

func methodColor(method string, colors fiber.Colors) string {
	switch method {
	case fiber.MethodGet:
		return colors.Cyan
	case fiber.MethodPost:
		return colors.Green
	case fiber.MethodPut:
		return colors.Yellow
	case fiber.MethodDelete:
		return colors.Red
	case fiber.MethodPatch:
		return colors.White
	case fiber.MethodHead:
		return colors.Magenta
	case fiber.MethodOptions:
		return colors.Blue
	default:
		return colors.Reset
	}
}

func NewLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		chainErr := c.Next()

		errHandler := c.App().ErrorHandler

		if chainErr != nil {
			logger.Error("%v", chainErr)
			if err := errHandler(c, chainErr); err != nil {
				_ = c.SendStatus(fiber.StatusInternalServerError)
			}
		}

		end := time.Now()

		if logger.DefaultConfig.Colors {
			logger.Info(
				"[%s]:%s  %s%s\x1b[0m  %s  %s%v\x1b[0m  %s%v\x1b[0m",
				c.IP(),
				c.Port(),
				methodColor(c.Method(), fiber.DefaultColors),
				c.Method(),
				c.Path(),
				statusColor(c.Response().StatusCode(), fiber.DefaultColors),
				c.Response().StatusCode(),
				"\x1b[90m",
				end.Sub(start),
			)
		} else {
			logger.Info(
				"[%s]:%s  %s  %s  %v  %v",
				c.IP(),
				c.Port(),
				c.Method(),
				c.Path(),
				c.Response().StatusCode(),
				end.Sub(start),
			)
		}

		return nil
	}
}
