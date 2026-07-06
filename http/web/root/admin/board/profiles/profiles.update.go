package profiles

import (
	"reflect"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type ProfilesUpdateForm struct {
	EnablePicture            bool     `form:"enable_picture"`
	PictureUploadFormats     []string `form:"upload_formats" validate:"required_if=EnablePicture true,dive,oneof=image/gif image/jpeg image/png image/webp"`
	PictureFormat            string   `form:"picture_format" validate:"required,oneof=webp png jpg"`
	PictureMaxSize           int64    `form:"picture_max_size" validate:"required,min=1"`
	PictureStorageProviderID string   `form:"picture_storage_provider_id" validate:"required_if=EnablePicture true,max=64"`
	PictureStoragePath       string   `form:"picture_storage_path" validate:"omitempty,max=255"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminBoardProfiles")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ProfilesUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	storages, err := r.Runtime.Config.Storages()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	valid := frm.PictureStorageProviderID == ""
	for _, storage := range storages {
		if storage.ID == frm.PictureStorageProviderID {
			valid = true
			break
		}
	}
	if !valid {
		req.Flash.SetError(errs.ErrInvalidStorageProvider)
		return req.RedirectToRoute(myRoute)
	}

	var settingProfiles *setting.Setting[setting.Profiles]
	settingProfiles, err = settingRepo.GetByID[setting.Profiles](
		r.Runtime.Repositories.Setting,
		"profiles",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	settingProfiles.JSONValue.EnablePicture = frm.EnablePicture
	settingProfiles.JSONValue.PictureUploadFormats = frm.PictureUploadFormats
	settingProfiles.JSONValue.PictureFormat = frm.PictureFormat
	settingProfiles.JSONValue.PictureMaxSize = frm.PictureMaxSize
	settingProfiles.JSONValue.PictureStorageProviderID = frm.PictureStorageProviderID
	settingProfiles.JSONValue.PictureStoragePath = frm.PictureStoragePath

	err = settingRepo.Update[setting.Profiles](
		r.Runtime.Repositories.Setting,
		settingProfiles,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
