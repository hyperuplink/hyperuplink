package board

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"sections": []string{
			"attachments",
			"categories",
			"forums",
			"profiles",
			"themes",
			"topics",
		},
	})
}
