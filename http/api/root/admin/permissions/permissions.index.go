package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
)

// @Summary		Show the groups and their permissions
// @Tags			admin
// @Produce		json
// @Success		200	{object}	logicpermissions.View
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/permissions [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	shown, err := logicpermissions.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(shown)
}
