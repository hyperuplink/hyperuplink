package settings

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func viewToggle(profile *setting.UserProfile, view string) *bool {
	switch view {
	case setting.UserProfileViewBanner:
		return &profile.ShowBanner
	case setting.UserProfileViewFooter:
		return &profile.ShowFooter
	case setting.UserProfileViewProfilePictures:
		return &profile.ShowProfilePictures
	}

	return nil
}

func (r *Route) View(c fiber.Ctx) (err error) {
	myRoute := route.For("AccountSettingsView")
	parentRoute := route.For("AccountSettings")
	req := request.New(r, c, myRoute,
		[]string{"base"}, parentRoute.AsURL()+"/index",
		parentRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var usr *user.User
	usr, err = req.GetUser()
	if ret, rerr := req.RedirectBackOnError(err); ret == true {
		return rerr
	}

	var settingUserProfile *setting.Setting[setting.UserProfile]
	settingUserProfile, err = settingRepo.GetOrCreateUserProfile(
		r.Runtime.Repositories.Setting,
		usr.ID,
	)
	if ret, rerr := req.RedirectBackOnError(err); ret == true {
		return rerr
	}

	toggle := viewToggle(&settingUserProfile.JSONValue, c.Params("view"))
	if toggle == nil {
		return req.RedirectBack()
	}
	*toggle = !*toggle

	err = settingRepo.Update[setting.UserProfile](
		r.Runtime.Repositories.Setting,
		settingUserProfile,
	)
	if ret, rerr := req.RedirectBackOnError(err); ret == true {
		return rerr
	}

	return req.RedirectBack()
}
