package keys

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
)

// @Summary		Revoke a key
// @Description	Deletes one of the signed-in user's keys. A key of another user is
// @Description	reported as not found.
// @Tags			account
// @Param			key	path	string	true	"The key UUID"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/keys/{key} [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("key")
	if ret {
		return rerr
	}

	err = logicapi.Destroy(r.Runtime, req.User().ID, &logicapi.DestroyInput{
		ID: id.String(),
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
