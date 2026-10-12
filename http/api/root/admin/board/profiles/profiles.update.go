package profiles

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicprofiles "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/profiles"
)

// @Summary		Update the profile settings
// @Tags			admin
// @Accept			json
// @Param			request	body	logicprofiles.UpdateInput	true	"The settings to store"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/profiles [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicprofiles.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	err = logicprofiles.Update(req.Context(), r.Runtime, req.ViewerID(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
