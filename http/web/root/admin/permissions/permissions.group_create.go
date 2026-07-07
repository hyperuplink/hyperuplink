package permissions

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/group"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

type GroupCreateForm struct {
	ID   string `form:"id" validate:"required,slug,min=1,max=32"`
	Name string `form:"name" validate:"required,min=1,max=32"`
}

func (r *Route) GroupCreate(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(GroupCreateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	grp := new(group.Group)
	grp.ID = frm.ID
	grp.Name = frm.Name

	_, err = r.Runtime.Repositories.Group.Create(grp)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
