package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
)

// @Summary		Revoke a group's permission in a category, or a category's default
// @Tags			admin
// @Accept			json
// @Param			request	body	logicpermissions.RemoveInput	true	"The category, and the group unless the category's default is meant"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/permissions/remove [post]
func (r *Route) Remove(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicpermissions.RemoveInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.Remove(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
