package permissions

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Apply(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminPermissions")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicpermissions.ApplyInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	err = logicpermissions.Apply(r.Runtime, in)
	if errors.Is(err, errs.ErrValidation) {
		req.Flash.SetError(err)
		return req.RedirectToRoute(myRoute)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
