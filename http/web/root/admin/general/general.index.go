package general

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminGeneral")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingGeneral *setting.Setting[setting.General]
	settingGeneral, err = settingRepo.GetByID[setting.General](
		r.Runtime.Repositories.Setting,
		"general",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_general", &settingGeneral.JSONValue)

	return req.Respond()
}
