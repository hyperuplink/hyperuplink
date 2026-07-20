package profile

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Show the account profile
// @Description	Returns the caller's own record together with the profile fields the
// @Description	board has enabled and the values the caller has filled in.
// @Tags			account
// @Produce		json
// @Success		200	{object}	object{user=user.Detail,profiles=setting.Profiles,user_profile=setting.UserProfile}
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/profile [get]
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
