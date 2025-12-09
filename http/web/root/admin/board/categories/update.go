package categories

import "github.com/gofiber/fiber/v3"

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}
