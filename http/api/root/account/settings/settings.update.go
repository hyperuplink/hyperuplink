package settings

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Update the account settings
// @Tags		account
// @Accept		json
// @Produce	json
// @Param		request	body		logicsettings.UpdateInput	true	"The settings to store"
// @Success	200		{object}	request.StatusResponse
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ValidationErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/account/settings [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicsettings.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicsettings.Update(r.Runtime, req.User, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
