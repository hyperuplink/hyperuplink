package users

import (
	"github.com/gofiber/fiber/v3"
	logicusers "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/users"
)

// @Summary		Lift a user's ban
// @Tags			admin
// @Param			user	path	string	true	"The user UUID"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/users/{user}/unban [post]
func (r *Route) Unban(c fiber.Ctx) (err error) {
	return r.action(c, logicusers.Unban)
}
