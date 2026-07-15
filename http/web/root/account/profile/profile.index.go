package profile

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
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

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var settingProfiles *setting.Setting[setting.Profiles]
	settingProfiles, err = settingRepo.GetByID[setting.Profiles](
		r.Runtime.Repositories.Setting,
		"profiles",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var settingUserProfile *setting.Setting[setting.UserProfile]
	settingUserProfile, err = settingRepo.GetOrCreateUserProfile(
		r.Runtime.Repositories.Setting,
		usr.ID,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)
	req.SetData("setting_profiles", &settingProfiles.JSONValue)
	req.SetData("setting_user_profile", &settingUserProfile.JSONValue)

	return req.Respond()
}
