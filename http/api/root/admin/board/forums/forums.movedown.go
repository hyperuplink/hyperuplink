package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	Move a forum down
// @Tags		admin
// @Produce	json
// @Param		id	path		string	true	"The forum identifier"
// @Success	200	{object}	request.StatusResponse
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Failure	404	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/board/forums/{id}/movedown [post]
func (r *Route) MoveDown(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	err = logicforums.MoveDown(r.Runtime, &logicforums.MoveInput{
		ID: c.Params("id"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
