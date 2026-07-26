package general

import (
	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

// @Summary	Show the general settings
// @Tags		admin
// @Produce	json
// @Success	200	{object}	object{system=setting.System,general=setting.General}
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/general [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	settingGeneral, err := settingRepo.GetByID[setting.General](
		gh.Repositories(r.Runtime).Setting,
		"general",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"system":  req.System,
		"general": settingGeneral.JSONValue,
	})
}
