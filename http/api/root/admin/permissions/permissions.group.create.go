package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Create a group
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		request	body		logicpermissions.GroupCreateInput	true	"The group to create"
// @Success	201		{object}	object{id=string}
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	409		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/permissions/group [post]
func (r *Route) GroupCreate(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicpermissions.GroupCreateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.GroupCreate(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondCreated(fiber.Map{"id": in.ID})
}
