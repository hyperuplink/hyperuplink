package settings

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) View(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	enabled, err := logicsettings.ToggleView(r.Runtime, req.User, c.Params("view"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"view":    c.Params("view"),
		"enabled": enabled,
	})
}
