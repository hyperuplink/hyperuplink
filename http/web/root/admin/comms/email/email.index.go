package email

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
	myRoute := route.For("AdminCommsEmail")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingCommsEmail *setting.Setting[setting.CommsEmail]
	settingCommsEmail, err = settingRepo.GetByID[setting.CommsEmail](
		r.Runtime.Repositories.Setting,
		"comms_email",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var targets config.Targets
	targets, err = r.Runtime.Config.Targets()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var emailTargets config.Targets
	var selectedTarget *config.Target
	for i := range targets {
		if !targets[i].Serves(config.TargetTypeEmail) {
			continue
		}
		emailTargets = append(emailTargets, targets[i])
		if targets[i].ID == settingCommsEmail.JSONValue.TargetID {
			selectedTarget = &targets[i]
		}
	}

	req.SetData("setting_comms_email", &settingCommsEmail.JSONValue)
	req.SetData("email_targets", emailTargets)
	req.SetData("selected_target", selectedTarget)

	return req.Respond()
}
