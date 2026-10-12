package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/topics"
)

// @Summary		Update the topic settings
// @Tags			admin
// @Accept			json
// @Param			request	body	logictopics.UpdateInput	true	"The settings to store"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/topics [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logictopics.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logictopics.Update(req.Context(), r.Runtime, req.ViewerID(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
