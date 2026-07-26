package api

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Create(c fiber.Ctx) (err error) {
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

	in := new(logicapi.CreateInput)

	if ok := req.ValidateForm(in, reflect.TypeOf(*in)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	userID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	var secret string
	_, secret, err = logicapi.Create(r.Runtime, userID, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.Flash.SetInfo(req.In.Ts("api_key_created") + " " + secret)
	return req.RedirectToRoute(myRoute)
}
