package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingTheme *setting.Setting[setting.Theme]
	settingTheme, err = settingRepo.GetByID[setting.Theme](
		r.Runtime.Repositories.Setting,
		"theme",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	themes, err := helpers.GetThemes(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	colorschemes, err := helpers.GetColorschemes(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var storageIDs []string
	for _, storage := range storages {
		storageIDs = append(storageIDs, storage.ID)
	}

	req.SetData("setting_theme", &settingTheme.JSONValue)
	req.SetData("themes", themes)
	req.SetData("colorschemes", colorschemes)
	req.SetData("storage_ids", storageIDs)

	return req.Respond()
}
