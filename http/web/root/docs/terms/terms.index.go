package terms

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("DocsTerms")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole,
		user.UserRole,
		user.ModRole,
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

	if !settingGeneral.JSONValue.EnableTerms {
		return c.SendStatus(fiber.StatusNotFound)
	}

	var html string
	html, err = r.Runtime.Markdown.Convert(settingGeneral.JSONValue.Terms)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("html", html)

	return req.Respond()
}
