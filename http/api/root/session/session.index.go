package session

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show the current session
// @Description	Returns the user the API key was minted by, which is as close as
// @Description	the key-authenticated API gets to a session.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{user=user.Detail}
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/session [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"user": req.User.AsDetail(),
	})
}
