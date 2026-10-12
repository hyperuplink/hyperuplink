package profiles

import (
	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

// @Summary		Show the profile settings
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/profiles [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	settingProfiles, err := settingRepo.GetByID[setting.Profiles](req.Context(), gh.Repositories(r.Runtime).Setting, "profiles")
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	storages, err := r.Runtime.Config().Storages()
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	storageIDs := []string{}
	for _, storage := range storages {
		storageIDs = append(storageIDs, storage.ID)
	}

	return req.Respond(fiber.Map{
		"profiles":    settingProfiles.JSONValue,
		"storage_ids": storageIDs,
	})
}
