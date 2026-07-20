package password

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpassword "xn--gckvb8fzb.com/hyperuplink/logic/root/account/password"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Change the account password
// @Description	The current password has to be sent along with the new one, and a
// @Description	user with two-factor authentication enabled has to add a current
// @Description	one-time code as well.
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			request	body		logicpassword.UpdateInput	true	"The current and the new password"
// @Success		200		{object}	request.StatusResponse
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/password [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(logicpassword.UpdateInput)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpassword.Update(r.Runtime, req.User, in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
