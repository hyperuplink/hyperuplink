package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
)

// @Summary		Delete a group
// @Tags			admin
// @Param			group	path	string	true	"The group id"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/permissions/groups/{group} [delete]
func (r *Route) GroupDestroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	err = logicpermissions.GroupDestroy(req.Context(), r.Runtime, &logicpermissions.GroupDestroyInput{
		ID: c.Params("group"),
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
