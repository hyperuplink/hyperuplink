package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logiccategories "xn--gckvb8fzb.com/hyperuplink/logic/root/categories"
)

// @Summary		Show a category
// @Description	One category with its forums, addressed by UUID or slug. A category
// @Description	the viewer may not read is reported as not found.
// @Tags			board
// @Produce		json
// @Param			category	path		string	true	"The category UUID or slug"
// @Success		200			{object}	object{category=view.Category}
// @Failure		401			{object}	problem.Problem
// @Failure		404			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/categories/{category} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	shown, err := logiccategories.Show(r.Runtime, c.Params("category"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	m := req.Views()

	return req.Respond(fiber.Map{
		"category": m.Category(shown.Category, m.Forums(*shown.Forums, req.Perms()), req.Perms()),
	})
}
