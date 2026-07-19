package profile

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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

	var view *logicprofile.View
	view, err = logicprofile.Show(r.Runtime, usr)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.SetData("user", usr)
	req.SetData("setting_profiles", view.Profiles)
	req.SetData("setting_user_profile", view.UserProfile)

	return req.Respond()
}
