package profile

import (
	"mime/multipart"
	"os"
	"reflect"

	"github.com/gofiber/fiber/v3"
	"github.com/lithammer/shortuuid/v4"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/request"
	"github.com/mrusme/hyperuplink/models/user"
)

type ProfileUpdateForm struct {
	ProfilePicture *multipart.FileHeader `form:"profile_picture" validate:"required"`
}

func (r *Route) Update(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountProfile")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Remove!
		user.UserRole,  // TODO: Remove!
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
	usr, err = r.getUser(req)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	profilePictureMultipartFile, err := frm.ProfilePicture.Open()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	profilePictureFile, profilePictureFileName, err := r.Runtime.Magick.Convert(
		profilePictureMultipartFile,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		profilePictureMultipartFile.Close()
		return rerr
	}

	profilePictureID := shortuuid.New()

	// TODO: Get provider ID from System
	// TODO: Get path from System
	// TODO: Get .webp from format configured in System
	err = r.Runtime.Storage.StoreFile(
		"profile-pictures",
		profilePictureFile,
		"profile-pictures/"+profilePictureID+".webp",
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

	err = r.Runtime.Repositories.User.Update(usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	// TODO: Delete oldProfilePicutreID

	return req.RedirectToRoute(myRoute)
}
