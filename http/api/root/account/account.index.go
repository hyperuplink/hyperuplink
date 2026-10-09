package account

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
)

// @Summary		Show the account
// @Description	The signed-in user's own record, with the email address, the
// @Description	timezone, the groups and the notification preference.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{user=view.Account}
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	return r.respondAccount(req)
}
