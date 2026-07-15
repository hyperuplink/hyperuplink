package xmpp

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/config"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type XMPPUpdateForm struct {
	TargetID string `form:"target_id" validate:"omitempty,max=64"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminCommsXmpp")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(XMPPUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	if frm.TargetID != "" {
		targets, terr := r.Runtime.Config.Targets()
		if ret, rerr := req.RespondOnError(terr); ret == true {
			return rerr
		}

		valid := false
		for _, target := range targets {
			if target.ID == frm.TargetID &&
				target.Serves(config.TargetTypeXMPP) {
				valid = true
				break
			}
		}
		if !valid {
			req.Flash.SetError(errs.ErrTargetIDNotFound)
			return req.RedirectToRoute(myRoute)
		}
	}

	var settingCommsXMPP *setting.Setting[setting.CommsXMPP]
	settingCommsXMPP, err = settingRepo.GetByID[setting.CommsXMPP](
		r.Runtime.Repositories.Setting,
		"comms_xmpp",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	before := settingCommsXMPP.JSONValue

	settingCommsXMPP.JSONValue.TargetID = frm.TargetID

	err = settingRepo.Update[setting.CommsXMPP](
		r.Runtime.Repositories.Setting,
		settingCommsXMPP,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if actorID, ok := req.Session.GetUserUUID(); ok {
		logicactivity.RecordAdminSettingsUpdate(r.Runtime, actorID,
			"comms_xmpp", before, settingCommsXMPP.JSONValue)
	}

	return req.RedirectToRoute(myRoute)
}
