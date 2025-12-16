package profile

import (
	"mime/multipart"
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

	profilePictureFile, err := frm.ProfilePicture.Open()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}
	defer profilePictureFile.Close()

	var usr *user.User
	usr, err = r.getUser(req)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	profilePictureID := shortuuid.New()

	// TODO: Get providerID from System config
	err = r.Runtime.Storage.StoreFile(
		"profile-pictures",
		profilePictureFile,
		"profile-pictures/"+profilePictureID+".png",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
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
