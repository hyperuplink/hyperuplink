package users

import (
	"github.com/gofiber/fiber/v3"
	logicusers "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/users"
)

func (r *Route) Destroy(c fiber.Ctx) (err error) {
	return r.action(c, logicusers.Delete)
}
