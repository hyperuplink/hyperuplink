package password

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicpassword "xn--gckvb8fzb.com/hyperuplink/logic/root/account/password"
)

// @Summary		Change the password
// @Description	The current password has to be sent with the new one and its
// @Description	repetition, and an account with two-factor authentication adds
// @Description	the current code. Keys stay valid.
// @Tags			account
// @Accept			json
// @Param			request	body	logicpassword.UpdateInput	true	"The current and the new password"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/password [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(logicpassword.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicpassword.Update(r.Runtime, req.User(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
