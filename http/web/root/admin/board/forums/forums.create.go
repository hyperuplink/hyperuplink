package forums

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardForums")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicforums.CreateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", in)

	_, err = logicforums.Create(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
