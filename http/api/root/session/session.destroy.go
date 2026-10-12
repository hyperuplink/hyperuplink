package session

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
)

// @Summary		Sign out
// @Description	Revokes the key the request was made with. The next request with
// @Description	that key is rejected.
// @Tags			session
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/session [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	err = logicapi.Destroy(req.Context(), r.Runtime, req.User().ID, &logicapi.DestroyInput{
		ID: req.Principal.Key.ID.String(),
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
