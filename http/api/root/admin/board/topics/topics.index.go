package topics

import (
	"github.com/gofiber/fiber/v3"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

// @Summary	Show the topic settings
// @Tags		admin
// @Produce	json
// @Success	200	{object}	object{topics=setting.Topics}
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/topics [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	settingTopics, err := settingRepo.GetByID[setting.Topics](
		gh.Repositories(r.Runtime).Setting,
		"topics",
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"topics": settingTopics.JSONValue,
	})
}
