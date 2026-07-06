package profile

import (
	"io"
	"mime/multipart"
	"os"
	"reflect"

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

type ProfileUpdateForm struct {
	ProfilePicture *multipart.FileHeader `form:"profile_picture" validate:""`
	SignatureText  string                `form:"signature_text" validate:"max=256"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountProfile")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	frm := new(ProfileUpdateForm)

	if ok := req.ValidateForm(frm, reflect.TypeOf(*frm)); !ok {
		return req.RedirectToRoute(myRoute)
	}

	r.Runtime.Debug("form", frm)

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	if frm.ProfilePicture != nil && frm.ProfilePicture.Filename != "" {
		var settingProfiles *setting.Setting[setting.Profiles]
		settingProfiles, err = settingRepo.GetByID[setting.Profiles](
			r.Runtime.Repositories.Setting,
			"profiles",
		)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}
		profiles := settingProfiles.JSONValue

		if profiles.EnablePicture {
			if frm.ProfilePicture.Size > profiles.GetPictureMaxSize() {
				req.Flash.SetError(errs.ErrPictureTooLarge)
				return req.RedirectToRoute(myRoute)
			}

			profilePictureMultipartFile, err := frm.ProfilePicture.Open()
			if ret, rerr := req.RespondOnError(err); ret == true {
				return rerr
			}

			mtype, err := mimetype.DetectReader(profilePictureMultipartFile)
			if ret, rerr := req.RespondOnError(err); ret == true {
				profilePictureMultipartFile.Close()
				return rerr
			}
			if _, err = profilePictureMultipartFile.Seek(0, io.SeekStart); err != nil {
				profilePictureMultipartFile.Close()
				return req.RespondError(err)
			}

			allowed := false
			for _, format := range profiles.PictureUploadFormats {
				if mtype.Is(format) {
					allowed = true
					break
				}
			}
			if !allowed {
				profilePictureMultipartFile.Close()
				req.Flash.SetError(errs.ErrPictureFormatNotAllowed)
				return req.RedirectToRoute(myRoute)
			}

			profilePictureFile, profilePictureFileName, err := r.Runtime.Magick.ConvertProfilePicture(
				profilePictureMultipartFile,
				profiles.PictureFormat,
			)
			if ret, rerr := req.RespondOnError(err); ret == true {
				profilePictureMultipartFile.Close()
				return rerr
			}

			profilePictureID := shortuuid.New()

			err = r.Runtime.Storage.StoreFile(
				profiles.PictureStorageProviderID,
				profilePictureFile,
				profiles.PictureStoragePath+"/"+profilePictureID+"."+profiles.PictureFormat,
			)
			if ret, rerr := req.RespondOnError(err); ret == true {
				return rerr
			}

			profilePictureFile.Close()
			if profilePictureFileName != "" {
				os.Remove(profilePictureFileName)
			}

			// oldProfilePicutreID := usr.ProfilePicture
			usr.ProfilePicture = profilePictureID
			// TODO: Delete oldProfilePicutreID
		}
	}

	usr.SignatureText = frm.SignatureText
	if usr.SignatureText == "" {
		usr.SignatureHTML = ""
	} else {
		usr.SignatureHTML, err = r.Runtime.Markdown.Convert(usr.SignatureText)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}
	}

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RedirectToRoute(myRoute)
}
