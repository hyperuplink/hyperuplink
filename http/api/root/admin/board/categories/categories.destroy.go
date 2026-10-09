package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
)

// @Summary		Delete a category
// @Description	The delete is soft and propagates to everything below.
// @Tags			admin
// @Param			category	path	string	true	"The category UUID"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/categories/{category} [delete]
func (r *Route) Destroy(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("category")
	if ret {
		return rerr
	}

	err = logiccategories.Destroy(r.Runtime, &logiccategories.DestroyInput{ID: id.String()})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
