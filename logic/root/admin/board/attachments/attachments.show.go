package attachments

import (
	"context"
	"strings"

	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Show(ctx context.Context, rt *runtime.Runtime) (view *View, err error) {
	settingAttachments, err := settingRepo.GetByID[setting.Attachments](ctx,
		gh.Repositories(rt).Setting,
		"attachments",
	)
	if err != nil {
		return nil, err
	}

	storages, err := rt.Config().Storages()
	if err != nil {
		return nil, err
	}

	view = new(View)
	view.Attachments = &settingAttachments.JSONValue
	view.UploadFormatOptions = setting.AttachmentUploadFormatOptions
	for _, storage := range storages {
		view.StorageIDs = append(view.StorageIDs, storage.ID)

		if storage.ID == settingAttachments.JSONValue.StorageProviderID &&
			strings.ToLower(storage.Type) == "s3" &&
			storage.S3.PublicDownload {
			view.StorageIsPublic = true
		}
	}

	return view, nil
}
