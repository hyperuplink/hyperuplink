package manual

import (
	"github.com/gofiber/fiber/v3"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	return r.render(c, "")
}
