package twofactor

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictwofactor "xn--gckvb8fzb.com/hyperuplink/logic/root/account/twofactor"
)

// @Summary		Disable two-factor authentication
// @Description	Turns two-factor authentication off; the current password confirms
// @Description	the request.
// @Tags			account
// @Accept			json
// @Param			request	body	twofactor.DisableBody	true	"The current password"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/twofactor [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(DisableBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logictwofactor.Disable(r.Runtime, req.User(), in.CurrentPassword)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
