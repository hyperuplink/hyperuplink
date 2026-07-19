package settings

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicsettings.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicsettings.Update(r.Runtime, req.User, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
