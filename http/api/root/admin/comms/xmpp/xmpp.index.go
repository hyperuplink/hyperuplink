package xmpp

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicxmpp "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/comms/xmpp"
)

// @Summary		Show the XMPP delivery settings
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/comms/xmpp [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	shown, err := logicxmpp.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	targetIDs := []string{}
	for _, target := range shown.XMPPTargets {
		targetIDs = append(targetIDs, target.ID)
	}

	return req.Respond(fiber.Map{
		"comms_xmpp": shown.CommsXMPP,
		"target_ids": targetIDs,
	})
}
