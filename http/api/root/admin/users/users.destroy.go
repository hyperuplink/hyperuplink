package users

import (
	"github.com/gofiber/fiber/v3"
	logicusers "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/users"
)

// @Summary		Delete a user
// @Description	The delete is soft; the user's posts stay.
// @Tags			admin
// @Param			user	path	string	true	"The user UUID"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/users/{user} [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	return r.action(c, logicusers.Delete)
}
