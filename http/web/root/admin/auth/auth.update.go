package auth

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type AuthUpdateForm struct {
	AddressType int `form:"address_type" validate:"oneof=0 1 2"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminAuth")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(AuthUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	var settingAuth *setting.Setting[setting.Auth]
	settingAuth, err = settingRepo.GetByID[setting.Auth](
		r.Runtime.Repositories.Setting,
		"auth",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	before := settingAuth.JSONValue

	settingAuth.JSONValue.AddressType = setting.AddressType(frm.AddressType)

	err = settingRepo.Update[setting.Auth](
		r.Runtime.Repositories.Setting,
		settingAuth,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if actorID, ok := req.Session.GetUserUUID(); ok {
		logicactivity.RecordAdminSettingsUpdate(r.Runtime, actorID,
			"auth", before, settingAuth.JSONValue)
	}

	return req.RedirectToRoute(myRoute)
}
