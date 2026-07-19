package profiles

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	settingProfiles, err := settingRepo.GetByID[setting.Profiles](
		r.Runtime.Repositories.Setting,
		"profiles",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var storageIDs []string
	for _, storage := range storages {
		storageIDs = append(storageIDs, storage.ID)
	}

	return req.Respond(fiber.Map{
		"profiles":    settingProfiles.JSONValue,
		"storage_ids": storageIDs,
	})
}
