package user

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Membership(c fiber.Ctx) (err error) {
	myRoute := route.For("User")
	req := request.New(r, c, myRoute,
		[]string{"base"}, "user/index",
		"")

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	username := c.Params("user")

	in := new(logicuser.MembershipInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute.Fill(map[string]string{"user": username}))
	}

	err = logicuser.UpdateMembership(r.Runtime, username, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute.Fill(map[string]string{"user": username}))
}
