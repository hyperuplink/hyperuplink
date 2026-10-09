package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

// @Summary		Show the theme settings
// @Tags			admin
// @Produce		json
// @Success		200	{object}	logicthemes.View
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	shown, err := logicthemes.Show(r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(shown)
}
