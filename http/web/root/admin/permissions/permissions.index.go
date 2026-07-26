package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicpermissions.View
	view, err = logicpermissions.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("default_level", view.DefaultLevel)
	req.SetData("groups", view.Groups)
	req.SetData("has_categories", view.HasCategories)

	return req.Respond()
}
