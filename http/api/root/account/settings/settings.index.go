package settings

import (
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/zlasd/tzloc"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	timezones := tzloc.GetLocationList()
	slices.Sort(timezones)

	return req.Respond(fiber.Map{
		"user":      req.User.AsDetail(),
		"timezones": timezones,
	})
}
