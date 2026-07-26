package auth

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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
		gh.Repositories(r.Runtime).Setting,
		"auth",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_auth", &settingAuth.JSONValue)

	var authProviders gh.AuthProviders
	_, err = r.Runtime.Config().Unmarshal("AuthProvider", &authProviders)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("auth_providers", authProviders)

	return req.Respond()
}
