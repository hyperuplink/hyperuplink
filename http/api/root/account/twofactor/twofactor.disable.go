package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Disable(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(TwofactorDisableInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logictwofactor.Disable(r.Runtime, req.User, in.CurrentPassword)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
