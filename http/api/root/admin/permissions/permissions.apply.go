package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
)

// @Summary		Set a permission level
// @Description	Without group_id and category_id this sets the board's default
// @Description	level; with category_id alone it sets that category's default for
// @Description	everyone, where the level "inherit" removes it again; with both it
// @Description	grants the group that level in the category.
// @Tags			admin
// @Accept			json
// @Param			request	body	logicpermissions.ApplyInput	true	"The level to set"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/permissions/apply [post]
func (r *Route) Apply(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicpermissions.ApplyInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.Apply(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
