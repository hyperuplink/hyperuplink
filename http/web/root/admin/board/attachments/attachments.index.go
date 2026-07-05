package attachments

import (
	"github.com/gofiber/fiber/v3"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardAttachments")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingAttachments *setting.Setting[setting.Attachments]
	settingAttachments, err = settingRepo.GetByID[setting.Attachments](
		r.Runtime.Repositories.Setting,
		"attachments",
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

	req.SetData("setting_attachments", &settingAttachments.JSONValue)
	req.SetData("storage_ids", storageIDs)
	req.SetData("upload_format_options", setting.AttachmentUploadFormatOptions)

	return req.Respond()
}
