package email

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicemail "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/comms/email"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

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

	targetIDs := []string{}
	for _, target := range view.EmailTargets {
		targetIDs = append(targetIDs, target.ID)
	}

	return req.Respond(fiber.Map{
		"comms_email": view.CommsEmail,
		"target_ids":  targetIDs,
	})
}
