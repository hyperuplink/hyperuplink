package themes

import (
	"io"
	"mime/multipart"
	"path"
	"reflect"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gofiber/fiber/v3"
	"github.com/lithammer/shortuuid/v4"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type ThemeBackgroundForm struct {
	CustomBackground *multipart.FileHeader `form:"custom_background" validate:""`
}

func (r *Route) BackgroundUpload(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ThemeBackgroundForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	if frm.CustomBackground == nil || frm.CustomBackground.Filename == "" {
		req.Flash.SetError(errs.ErrBackgroundUploadFailed)
		return req.RedirectToRoute(myRoute)
	}

	var settingTheme *setting.Setting[setting.Theme]
	settingTheme, err = settingRepo.GetByID[setting.Theme](
		r.Runtime.Repositories.Setting,
		"theme",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	theme := settingTheme.JSONValue

	if theme.ThemeStorageProviderID == "" {
		req.Flash.SetError(errs.ErrThemeStorageNotConfigured)
		return req.RedirectToRoute(myRoute)
	}

	if theme.CustomBackground != "" {
		req.Flash.SetError(errs.ErrBackgroundAlreadySet)
		return req.RedirectToRoute(myRoute)
	}

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	valid := false
	for _, storage := range storages {
		if storage.ID == theme.ThemeStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		req.Flash.SetError(errs.ErrInvalidStorageProvider)
		return req.RedirectToRoute(myRoute)
	}

	backgroundFile, err := frm.CustomBackground.Open()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	defer backgroundFile.Close()

	mtype, err := mimetype.DetectReader(backgroundFile)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if _, err = backgroundFile.Seek(0, io.SeekStart); err != nil {
		return req.RespondError(err)
	}

	if !strings.HasPrefix(mtype.String(), "image/") {
		req.Flash.SetError(errs.ErrBackgroundFormatNotAllowed)
		return req.RedirectToRoute(myRoute)
	}

	backgroundFilename := shortuuid.New() + mtype.Extension()

	err = r.Runtime.Storage.StoreFile(
		theme.ThemeStorageProviderID,
		backgroundFile,
		path.Join(theme.ThemeStoragePath, backgroundFilename),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingTheme.JSONValue.CustomBackground = backgroundFilename

	err = settingRepo.Update[setting.Theme](
		r.Runtime.Repositories.Setting,
		settingTheme,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}

func (r *Route) BackgroundRemove(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var settingTheme *setting.Setting[setting.Theme]
	settingTheme, err = settingRepo.GetByID[setting.Theme](
		r.Runtime.Repositories.Setting,
		"theme",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	theme := settingTheme.JSONValue

	if theme.CustomBackground != "" && theme.ThemeStorageProviderID != "" {
		if derr := r.Runtime.Storage.DeleteFile(
			theme.ThemeStorageProviderID,
			path.Join(theme.ThemeStoragePath, theme.CustomBackground),
		); derr != nil {
			r.Runtime.Warn("error", derr)
		}
	}

	settingTheme.JSONValue.CustomBackground = ""

	err = settingRepo.Update[setting.Theme](
		r.Runtime.Repositories.Setting,
		settingTheme,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
