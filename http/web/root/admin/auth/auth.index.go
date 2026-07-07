package auth

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminAuth")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingAuth *setting.Setting[setting.Auth]
	settingAuth, err = settingRepo.GetByID[setting.Auth](
		r.Runtime.Repositories.Setting,
		"auth",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_auth", &settingAuth.JSONValue)

	var authProviders config.AuthProviders
	authProviders, err = r.Runtime.Config.AuthProviders()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("auth_providers", authProviders)

	return req.Respond()
}
