package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
)

// @Summary		Update the theme settings
// @Tags			admin
// @Accept			json
// @Param			request	body	logicthemes.UpdateInput	true	"The theme settings to store"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/themes [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicthemes.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicthemes.Update(r.Runtime, req.ViewerID(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
