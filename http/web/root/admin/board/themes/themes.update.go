package themes

import (
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type ThemesUpdateForm struct {
	Theme                  string `form:"theme" validate:"required"`
	Colorscheme            string `form:"colorscheme" validate:"required"`
	ThemeStorageProviderID string `form:"theme_storage_provider_id" validate:"omitempty,max=64"`
	ThemeStoragePath       string `form:"theme_storage_path" validate:"omitempty,max=255"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ThemesUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	themes, err := helpers.GetThemes(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if !slices.Contains(themes, frm.Theme) {
		req.Flash.SetError(errs.ErrInvalidTheme)
		return req.RedirectToRoute(myRoute)
	}

	colorschemes, err := helpers.GetColorschemes(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if !slices.Contains(colorschemes, frm.Colorscheme) {
		req.Flash.SetError(errs.ErrInvalidColorscheme)
		return req.RedirectToRoute(myRoute)
	}

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	valid := frm.ThemeStorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == frm.ThemeStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		req.Flash.SetError(errs.ErrInvalidStorageProvider)
		return req.RedirectToRoute(myRoute)
	}

	var settingTheme *setting.Setting[setting.Theme]
	settingTheme, err = settingRepo.GetByID[setting.Theme](
		r.Runtime.Repositories.Setting,
		"theme",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	before := settingTheme.JSONValue

	settingTheme.JSONValue.Theme = frm.Theme
	settingTheme.JSONValue.Colorscheme = frm.Colorscheme
	settingTheme.JSONValue.ThemeStorageProviderID = frm.ThemeStorageProviderID
	settingTheme.JSONValue.ThemeStoragePath = frm.ThemeStoragePath

	err = settingRepo.Update[setting.Theme](
		r.Runtime.Repositories.Setting,
		settingTheme,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if actorID, ok := req.Session.GetUserUUID(); ok {
		logicactivity.RecordAdminSettingsUpdate(r.Runtime, actorID,
			"theme", before, settingTheme.JSONValue)
	}

	return req.RedirectToRoute(myRoute)
}
