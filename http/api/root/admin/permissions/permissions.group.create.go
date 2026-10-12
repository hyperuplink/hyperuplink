package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
)

// @Summary		Create a group
// @Tags			admin
// @Accept			json
// @Produce		json
// @Param			request	body		logicpermissions.GroupCreateInput	true	"The group to create"
// @Success		201		{object}	object{id=string}
// @Header			201		{string}	Location	"The path of the new group"
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		409		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/permissions/groups [post]
func (r *Route) GroupCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicpermissions.GroupCreateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.GroupCreate(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	c.Set(fiber.HeaderLocation, view.APIPath("APIV1AdminPermissionsGroup", map[string]string{
		"group": in.ID,
	}))

	return req.RespondCreated(fiber.Map{"id": in.ID})
}
