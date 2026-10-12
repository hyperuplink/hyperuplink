package keys

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
)

// @Summary		List the keys
// @Description	Every API key of the signed-in user, newest first, without the
// @Description	secrets, which the board doesn't keep.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{keys=[]view.Key}
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/keys [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	keys, err := logicapi.List(req.Context(), r.Runtime, req.User().ID)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"keys": view.Keys(keys)})
}
