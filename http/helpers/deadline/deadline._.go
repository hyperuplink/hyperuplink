package deadline

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v3"
)

func New(timeout time.Duration) fiber.Handler {
	return func(c fiber.Ctx) error {
		if timeout <= 0 {
			return c.Next()
		}

		parent := c.Context()
		ctx, cancel := context.WithTimeout(parent, timeout)
		defer cancel()

		c.SetContext(ctx)
		err := c.Next()
		c.SetContext(parent)

		return err
	}
}
