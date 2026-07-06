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
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/setting"
	"github.com/mrusme/hyperuplink/models/user"
	settingRepo "github.com/mrusme/hyperuplink/services/repositories/setting"
)

type ThemeBannerForm struct {
	CustomBanner *multipart.FileHeader `form:"custom_banner" validate:""`
}

func (r *Route) BannerUpload(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ThemeBannerForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	if frm.CustomBanner == nil || frm.CustomBanner.Filename == "" {
		req.Flash.SetError(errs.ErrBannerUploadFailed)
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

	if theme.CustomBanner != "" {
		req.Flash.SetError(errs.ErrBannerAlreadySet)
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

	bannerFile, err := frm.CustomBanner.Open()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	defer bannerFile.Close()

	mtype, err := mimetype.DetectReader(bannerFile)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	if _, err = bannerFile.Seek(0, io.SeekStart); err != nil {
		return req.RespondError(err)
	}

	if !strings.HasPrefix(mtype.String(), "image/") {
		req.Flash.SetError(errs.ErrBannerFormatNotAllowed)
		return req.RedirectToRoute(myRoute)
	}

	bannerFilename := shortuuid.New() + mtype.Extension()

	err = r.Runtime.Storage.StoreFile(
		theme.ThemeStorageProviderID,
		bannerFile,
		path.Join(theme.ThemeStoragePath, bannerFilename),
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingTheme.JSONValue.CustomBanner = bannerFilename

	err = settingRepo.Update[setting.Theme](
		r.Runtime.Repositories.Setting,
		settingTheme,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}

func (r *Route) BannerRemove(c fiber.Ctx) (err error) {
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

	if theme.CustomBanner != "" && theme.ThemeStorageProviderID != "" {
		if derr := r.Runtime.Storage.DeleteFile(
			theme.ThemeStorageProviderID,
			path.Join(theme.ThemeStoragePath, theme.CustomBanner),
		); derr != nil {
			r.Runtime.Warn("error", derr)
		}
	}

	settingTheme.JSONValue.CustomBanner = ""

	err = settingRepo.Update[setting.Theme](
		r.Runtime.Repositories.Setting,
		settingTheme,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
