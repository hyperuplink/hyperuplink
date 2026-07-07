package permissions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type GroupDestroyForm struct {
	ID string `form:"id" validate:"required,slug,min=1,max=32"`
}

func (r *Route) GroupDestroy(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(GroupDestroyForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	grp := new(group.Group)
	grp.ID = frm.ID

	err = r.Runtime.Repositories.Group.Delete(grp)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
