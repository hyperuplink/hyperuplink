package profiles

import (
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	EnablePicture            bool     `json:"enable_picture" form:"enable_picture"`
	PictureUploadFormats     []string `json:"upload_formats" form:"upload_formats" validate:"required_if=EnablePicture true,dive,oneof=image/gif image/jpeg image/png image/webp"`
	PictureFormat            string   `json:"picture_format" form:"picture_format" validate:"required,oneof=webp png jpg"`
	PictureMaxSize           int64    `json:"picture_max_size" form:"picture_max_size" validate:"required,min=1"`
	PictureStorageProviderID string   `json:"picture_storage_provider_id" form:"picture_storage_provider_id" validate:"required_if=EnablePicture true,max=64"`
	PictureStoragePath       string   `json:"picture_storage_path" form:"picture_storage_path" validate:"omitempty,max=255"`
}

func Update(
	rt *runtime.Runtime,
	actorID uuid.NullUUID,
	in *UpdateInput,
) (err error) {
	storages, err := rt.Config().Storages()
	if err != nil {
		return err
	}

	valid := in.PictureStorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == in.PictureStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		return errs.ErrInvalidStorageProvider
	}

	settingProfiles, err := settingRepo.GetByID[setting.Profiles](
		gh.Repositories(rt).Setting,
		"profiles",
	)
	if err != nil {
		return err
	}

	before := settingProfiles.JSONValue

	settingProfiles.JSONValue.EnablePicture = in.EnablePicture
	settingProfiles.JSONValue.PictureUploadFormats = in.PictureUploadFormats
	settingProfiles.JSONValue.PictureFormat = in.PictureFormat
	settingProfiles.JSONValue.PictureMaxSize = in.PictureMaxSize
	settingProfiles.JSONValue.PictureStorageProviderID = in.PictureStorageProviderID
	settingProfiles.JSONValue.PictureStoragePath = in.PictureStoragePath

	if err = settingRepo.Update[setting.Profiles](
		gh.Repositories(rt).Setting,
		settingProfiles,
	); err != nil {
		return err
	}

	if actorID.Valid {
		logicactivity.RecordAdminSettingsUpdate(rt, actorID.UUID,
			"profiles", before, settingProfiles.JSONValue)
	}

	return nil
}
