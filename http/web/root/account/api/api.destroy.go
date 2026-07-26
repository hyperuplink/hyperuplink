package api

import (
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Destroy(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountAPI")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicapi.DestroyInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	userID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	err = logicapi.Destroy(r.Runtime, userID, in)
	if errors.Is(err, errs.ErrNoRows) {
		return req.RedirectToRoute(myRoute)
	}
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
