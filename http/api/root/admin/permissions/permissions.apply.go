package permissions

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpermissions "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Grant a group a permission on a category
// @Tags		admin
// @Accept		json
// @Produce	json
// @Param		request	body		logicpermissions.ApplyInput	true	"The group, the category and the level to grant"
// @Success	200		{object}	request.StatusResponse
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/permissions/apply [post]
func (r *Route) Apply(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicpermissions.ApplyInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpermissions.Apply(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
