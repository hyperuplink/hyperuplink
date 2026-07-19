package profile

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var view *logicprofile.View
	view, err = logicprofile.Show(r.Runtime, req.User)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"user":         req.User.AsDetail(),
		"profiles":     view.Profiles,
		"user_profile": view.UserProfile,
	})
}
