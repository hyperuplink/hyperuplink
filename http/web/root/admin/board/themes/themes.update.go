package themes

import (
	"errors"
	"reflect"
	"slices"

	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/helpers"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

type ThemesUpdateForm struct {
	Theme       string `form:"theme" validate:"required"`
	Colorscheme string `form:"colorscheme" validate:"required"`
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
		req.Flash.SetError(errors.New("invalid_theme"))
		return req.RedirectToRoute(myRoute)
	}

	colorschemes, err := helpers.GetColorschemes(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if !slices.Contains(colorschemes, frm.Colorscheme) {
		req.Flash.SetError(errors.New("invalid_colorscheme"))
		return req.RedirectToRoute(myRoute)
	}

	var settingSystem *setting.Setting[setting.System]
	settingSystem, err = settingRepo.GetByID[setting.System](
		r.Runtime.Repositories.Setting,
		"system",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingSystem.JSONValue.Theme = frm.Theme
	settingSystem.JSONValue.Colorscheme = frm.Colorscheme

	err = settingRepo.Update[setting.System](
		r.Runtime.Repositories.Setting,
		settingSystem,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
