package password

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpassword "xn--gckvb8fzb.com/hyperuplink/logic/root/account/password"
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

	in := new(logicpassword.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpassword.Update(r.Runtime, req.User, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
