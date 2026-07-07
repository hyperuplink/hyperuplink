package attachments

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
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
	storageIsPublic := false
	for _, storage := range storages {
		storageIDs = append(storageIDs, storage.ID)

		if storage.ID == settingAttachments.JSONValue.StorageProviderID &&
			strings.ToLower(storage.Type) == "s3" &&
			storage.S3.PublicDownload {
			storageIsPublic = true
		}
	}

	req.SetData("setting_attachments", &settingAttachments.JSONValue)
	req.SetData("storage_ids", storageIDs)
	req.SetData("storage_is_public", storageIsPublic)
	req.SetData("upload_format_options", setting.AttachmentUploadFormatOptions)

	return req.Respond()
}
