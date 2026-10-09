package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
)

// @Summary		Move a forum down
// @Tags			admin
// @Param			forum	path	string	true	"The forum UUID"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/forums/{forum}/movedown [post]
func (r *Route) MoveDown(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("forum")
	if ret {
		return rerr
	}

	err = logicforums.MoveDown(r.Runtime, &logicforums.MoveInput{ID: id.String()})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
