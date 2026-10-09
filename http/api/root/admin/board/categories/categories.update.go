package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
)

// @Summary		Update a category
// @Tags			admin
// @Accept			json
// @Param			category	path	string	true	"The category UUID"
// @Param			request	body	logiccategories.UpdateInput	true	"The fields to store"
// @Success		204
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Failure		404	{object}	problem.Problem
// @Failure		422	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/categories/{category} [put]
func (r *Route) Update(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	id, ret, rerr := req.ParamUUID("category")
	if ret {
		return rerr
	}

	in := new(logiccategories.UpdateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}
	in.ID = id.String()

	err = logiccategories.Update(r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
