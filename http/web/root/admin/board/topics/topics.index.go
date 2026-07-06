package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardTopics")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingTopics *setting.Setting[setting.Topics]
	settingTopics, err = settingRepo.GetByID[setting.Topics](
		r.Runtime.Repositories.Setting,
		"topics",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("topics", &settingTopics.JSONValue)

	return req.Respond()
}
