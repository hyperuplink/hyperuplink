package profile

import (
	"io"
	"mime/multipart"
	"os"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type UpdateInput struct {
	ProfilePicture *multipart.FileHeader `json:"-" form:"profile_picture" validate:""`
	SignatureText  string                `json:"signature_text" form:"signature_text" validate:"max=256"`
	NotifyOnReply  bool                  `json:"notify_on_reply" form:"notify_on_reply"`
}

type View struct {
	Profiles    *setting.Profiles    `json:"profiles"`
	UserProfile *setting.UserProfile `json:"user_profile"`
}

func StorePicture(
	rt *runtime.Runtime,
	usr *user.User,
	fh *multipart.FileHeader,
) (err error) {
	settingProfiles, err := settingRepo.GetByID[setting.Profiles](
		gh.Repositories(rt).Setting,
		"profiles",
	)
	if err != nil {
		return err
	}
	profiles := settingProfiles.JSONValue

	if !profiles.EnablePicture {
		return nil
	}

	if fh.Size > profiles.GetPictureMaxSize() {
		return errs.ErrPictureTooLarge
	}

	pictureFile, err := fh.Open()
	if err != nil {
		return err
	}
	defer pictureFile.Close()

	mtype, err := mimetype.DetectReader(pictureFile)
	if err != nil {
		return err
	}
	if _, err = pictureFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	allowed := false
	for _, format := range profiles.PictureUploadFormats {
		if mtype.Is(format) {
			allowed = true
			break
		}
	}
	if !allowed {
		return errs.ErrPictureFormatNotAllowed
	}

	converted, convertedName, err := gh.Magick(rt).ConvertProfilePicture(
		pictureFile,
		profiles.PictureFormat,
	)
	if err != nil {
		return err
	}

	pictureID := shortuuid.New()

	if err = rt.Storage().StoreFile(
		profiles.PictureStorageProviderID,
		converted,
		profiles.PictureStoragePath+"/"+pictureID+"."+profiles.PictureFormat,
	); err != nil {
		return err
	}

	converted.Close()
	if convertedName != "" {
		os.Remove(convertedName)
	}

	// oldProfilePicutreID := usr.ProfilePicture
	usr.ProfilePicture = pictureID
	// TODO: Delete oldProfilePicutreID

	return nil
}
