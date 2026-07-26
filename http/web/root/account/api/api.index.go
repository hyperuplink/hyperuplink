package api

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
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

	userID, ok := req.Session.GetUserUUID()
	if !ok {
		return req.RedirectToRoot()
	}

	var keys *[]apikey.APIKey
	keys, err = logicapi.List(r.Runtime, userID)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("apikeys", keys)

	return req.Respond()
}
