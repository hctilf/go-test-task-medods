package middleware

import (
	"time"

	"github.com/gofiber/fiber/v2"

	app "github.com/hctilf/go-test-task-medods/internal/app"
)

func Logger(app *app.Application) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		err := c.Next()

		app.Logger.Infow(
			"request",
			"from", c.IP(),
			"method", c.Method(),
			"path", c.Path(),
			"status", c.Response().StatusCode(),
			"latency", time.Since(start),
			"error", err,
			"userAgent", c.Get("User-Agent"),
			"host", c.Get("Host"),
			"referer", c.Get("Referer"),
			"requestId", c.Response().Header.Peek("X-Request-Id"),
		)

		return err
	}
}
