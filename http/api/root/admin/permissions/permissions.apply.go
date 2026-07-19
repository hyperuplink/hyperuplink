package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Apply(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicpermissions.ApplyInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.Apply(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
