package profile

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
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

	in := new(logicprofile.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicprofile.Update(r.Runtime, req.User, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
