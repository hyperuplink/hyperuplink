package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Delete a forum
// @Description	The delete is soft and propagates down to the topics and replies the
// @Description	forum holds.
// @Tags			admin
// @Produce		json
// @Param			id	path		string	true	"The forum identifier"
// @Success		200	{object}	request.StatusResponse
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Failure		404	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/admin/board/forums/{id} [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	err = logicforums.Destroy(r.Runtime, &logicforums.DestroyInput{
		ID: c.Params("id"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
