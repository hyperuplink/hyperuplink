package xmpp

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicxmpp "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/comms/xmpp"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Show the XMPP settings
// @Tags		admin
// @Produce	json
// @Success	200	{object}	object{comms_xmpp=setting.CommsXMPP,target_ids=[]string}
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/comms/xmpp [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

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

	targetIDs := []string{}
	for _, target := range view.XMPPTargets {
		targetIDs = append(targetIDs, target.ID)
	}

	return req.Respond(fiber.Map{
		"comms_xmpp": view.CommsXMPP,
		"target_ids": targetIDs,
	})
}
