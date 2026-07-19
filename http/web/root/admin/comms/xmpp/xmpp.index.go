package xmpp

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicxmpp "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/comms/xmpp"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var view *logicxmpp.View
	view, err = logicxmpp.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("setting_comms_xmpp", view.CommsXMPP)
	req.SetData("xmpp_targets", view.XMPPTargets)
	req.SetData("selected_target", view.SelectedTarget)

	return req.Respond()
}
