package categories

import "github.com/gofiber/fiber/v3"

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}
