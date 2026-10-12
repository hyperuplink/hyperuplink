package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/categories"
)

// @Summary		Create a category
// @Tags			admin
// @Accept			json
// @Produce		json
// @Param			request	body		logiccategories.CreateInput	true	"The category to create"
// @Success		201		{object}	object{id=string}
// @Header			201		{string}	Location	"The path of the new category"
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		409		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/categories [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logiccategories.CreateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	id, err := logiccategories.Create(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	c.Set(fiber.HeaderLocation, view.APIPath("APIV1AdminBoardCategory", map[string]string{
		"category": id.String(),
	}))

	return req.RespondCreated(fiber.Map{"id": id})
}
