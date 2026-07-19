package attachments

import (
	"slices"

	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	for _, format := range in.UploadFormats {
		if !slices.Contains(setting.AttachmentUploadFormatOptions, format) {
			return errs.ErrInvalidUploadFormat
		}
	}

	storages, err := rt.Config.Storages()
	if err != nil {
		return err
	}

	valid := in.StorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == in.StorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		return errs.ErrInvalidStorageProvider
	}

	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		rt.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return err
	}

	before := settingAttachments.JSONValue

	settingAttachments.JSONValue.EnableAttachments = in.EnableAttachments
	settingAttachments.JSONValue.UploadFormats = in.UploadFormats
	settingAttachments.JSONValue.MaxSize = in.MaxSize
	settingAttachments.JSONValue.StorageProviderID = in.StorageProviderID
	settingAttachments.JSONValue.StoragePath = in.StoragePath
	settingAttachments.JSONValue.OnUploadHook = in.OnUploadHook
	settingAttachments.JSONValue.InlineImageDisplay = in.InlineImageDisplay

	if err = settingRepo.Update[setting.Attachments](
		rt.Repositories.Setting,
		settingAttachments,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"attachments", before, settingAttachments.JSONValue)
	}

	return nil
}
