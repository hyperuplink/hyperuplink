package user

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Membership(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicuser.MembershipInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicuser.UpdateMembership(r.Runtime, c.Params("user"), in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
