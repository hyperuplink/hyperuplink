package xmpp

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
	myRoute := route.For("AdminCommsXmpp")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingCommsXMPP *setting.Setting[setting.CommsXMPP]
	settingCommsXMPP, err = settingRepo.GetByID[setting.CommsXMPP](
		r.Runtime.Repositories.Setting,
		"comms_xmpp",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var targets config.Targets
	targets, err = r.Runtime.Config.Targets()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var xmppTargets config.Targets
	var selectedTarget *config.Target
	for i := range targets {
		if targets[i].Type != "xmpp" {
			continue
		}
		xmppTargets = append(xmppTargets, targets[i])
		if targets[i].ID == settingCommsXMPP.JSONValue.TargetID {
			selectedTarget = &targets[i]
		}
	}

	req.SetData("setting_comms_xmpp", &settingCommsXMPP.JSONValue)
	req.SetData("xmpp_targets", xmppTargets)
	req.SetData("selected_target", selectedTarget)

	return req.Respond()
}
