package keys

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicapi "xn--gckvb8fzb.com/hyperuplink/logic/root/account/api"
)

// @Summary		Issue a key
// @Description	Creates a named API key for the signed-in user and returns its
// @Description	secret once, in token.
// @Tags			account
// @Accept			json
// @Produce		json
// @Param			request	body		logicapi.CreateInput	true	"The name of the key"
// @Success		201		{object}	object{token=string,key=view.Key}
// @Failure		401		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/keys [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(logicapi.CreateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	key, secret, err := logicapi.Create(r.Runtime, req.User().ID, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"token": secret,
		"key":   view.KeyFrom(key),
	})
}
