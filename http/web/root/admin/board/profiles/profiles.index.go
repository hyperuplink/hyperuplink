package profiles

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardProfiles")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingProfiles *setting.Setting[setting.Profiles]
	settingProfiles, err = settingRepo.GetByID[setting.Profiles](
		gh.Repositories(r.Runtime).Setting,
		"profiles",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	storages, err := r.Runtime.Config().Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var storageIDs []string
	for _, storage := range storages {
		storageIDs = append(storageIDs, storage.ID)
	}

	req.SetData("setting_profiles", &settingProfiles.JSONValue)
	req.SetData("storage_ids", storageIDs)
	req.SetData("upload_format_options", setting.PictureUploadFormatOptions)

	return req.Respond()
}
