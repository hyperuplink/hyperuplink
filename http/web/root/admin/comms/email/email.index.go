package email

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicemail "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/comms/email"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var view *logicemail.View
	view, err = logicemail.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_comms_email", view.CommsEmail)
	req.SetData("email_targets", view.EmailTargets)
	req.SetData("selected_target", view.SelectedTarget)

	return req.Respond()
}
