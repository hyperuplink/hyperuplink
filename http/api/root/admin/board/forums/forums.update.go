package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
)

// @Summary		Update a forum
// @Tags			admin
// @Accept			json
// @Param			forum	path	string	true	"The forum UUID"
// @Param			request	body	logicforums.UpdateInput	true	"The fields to store"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/forums/{forum} [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("forum")
	if ret {
		return rerr
	}

	in := new(logicforums.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}
	in.ID = id.String()

	err = logicforums.Update(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
