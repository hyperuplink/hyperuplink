package settings

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicsettings "xn--gckvb8fzb.com/hyperuplink/logic/root/account/settings"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Toggle a collapsible view
// @Description	Flips whether the named section of the interface starts out open or
// @Description	closed for the caller and answers with the state it now has.
// @Tags			account
// @Produce		json
// @Param			view	path		string	true	"The name of the view to toggle"
// @Success		200		{object}	object{view=string,enabled=boolean}
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/settings/view/{view} [post]
func (r *Route) View(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	enabled, err := logicsettings.ToggleView(r.Runtime, req.User, c.Params("view"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"view":    c.Params("view"),
		"enabled": enabled,
	})
}
