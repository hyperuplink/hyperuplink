package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Delete a category
// @Description	The delete is soft and propagates down to the forums, topics and
// @Description	replies the category holds.
// @Tags			admin
// @Produce		json
// @Param			id	path		string	true	"The category identifier"
// @Success		200	{object}	request.StatusResponse
// @Failure		401	{object}	request.ErrorResponse
// @Failure		403	{object}	request.ErrorResponse
// @Failure		404	{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/admin/board/categories/{id} [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	err = logiccategories.Destroy(r.Runtime, &logiccategories.DestroyInput{
		ID: c.Params("id"),
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
