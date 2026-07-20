package users

import (
	"github.com/gofiber/fiber/v3"
	logicusers "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/users"
)

// @Summary	Lift a user's ban
// @Tags		admin
// @Produce	json
// @Param		id	path		string	true	"The user identifier"
// @Success	200	{object}	request.StatusResponse
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Failure	404	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/users/{id}/unban [post]
func (r *Route) Unban(c fiber.Ctx) (err error) {
	return r.action(c, logicusers.Unban)
}
