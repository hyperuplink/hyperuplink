package themes

import (
	"io"
	"mime/multipart"
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func UploadImage(
	rt *runtime.Runtime,
	kind ImageKind,
	fh *multipart.FileHeader,
) (err error) {
	if fh == nil || fh.Filename == "" {
		return kind.errUploadFailed
	}

	settingTheme, err := settingRepo.GetByID[setting.Theme](
		gh.Repositories(rt).Setting,
		"theme",
	)
	if err != nil {
		return err
	}
	theme := settingTheme.JSONValue

	if theme.ThemeStorageProviderID == "" {
		return errs.ErrThemeStorageNotConfigured
	}

	if *kind.field(&theme) != "" {
		return kind.errAlreadySet
	}

	storages, err := rt.Config().Storages()
	if err != nil {
		return err
	}

	valid := false
	for _, storage := range storages {
		if storage.ID == theme.ThemeStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		return errs.ErrInvalidStorageProvider
	}

	imageFile, err := fh.Open()
	if err != nil {
		return err
	}
	defer imageFile.Close()

	mtype, err := mimetype.DetectReader(imageFile)
	if err != nil {
		return err
	}
	if _, err = imageFile.Seek(0, io.SeekStart); err != nil {
		return err
	}

	if !strings.HasPrefix(mtype.String(), "image/") {
		return kind.errFormatInvalid
	}

	imageFilename := shortuuid.New() + mtype.Extension()

	if err = rt.Storage().StoreFile(
		theme.ThemeStorageProviderID,
		imageFile,
		path.Join(theme.ThemeStoragePath, imageFilename),
	); err != nil {
		return err
	}

	*kind.field(&settingTheme.JSONValue) = imageFilename

	return settingRepo.Update[setting.Theme](
		gh.Repositories(rt).Setting,
		settingTheme,
	)
}

func RemoveImage(
	rt *runtime.Runtime,
	kind ImageKind,
) (err error) {
	settingTheme, err := settingRepo.GetByID[setting.Theme](
		gh.Repositories(rt).Setting,
		"theme",
	)
	if err != nil {
		return err
	}
	theme := settingTheme.JSONValue

	filename := *kind.field(&theme)
	if filename != "" && theme.ThemeStorageProviderID != "" {
		if derr := rt.Storage().DeleteFile(
			theme.ThemeStorageProviderID,
			path.Join(theme.ThemeStoragePath, filename),
		); derr != nil {
			rt.Warn("error", derr)
		}
	}

	*kind.field(&settingTheme.JSONValue) = ""

	return settingRepo.Update[setting.Theme](
		gh.Repositories(rt).Setting,
		settingTheme,
	)
}
