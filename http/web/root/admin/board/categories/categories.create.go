package categories

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Create(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardCategories")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logiccategories.CreateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", in)

	_, err = logiccategories.Create(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
