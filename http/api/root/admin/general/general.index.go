package general

import (
	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

// @Summary		Show the general settings
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/general [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	settingGeneral, err := settingRepo.GetByID[setting.General](req.Context(), gh.Repositories(r.Runtime).Setting, "general")
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"system":  req.System,
		"general": settingGeneral.JSONValue,
	})
}
