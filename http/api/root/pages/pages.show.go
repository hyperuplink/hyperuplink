package pages

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicdocs "xn--gckvb8fzb.com/hyperuplink/logic/root/docs"
)

// @Summary		Show a board page
// @Description	One of the texts the administrator can enable: about, contact,
// @Description	privacy or terms. A page that is disabled is reported as not found.
// @Tags			board
// @Produce		json
// @Param			page	path		string	true	"about, contact, privacy or terms"
// @Success		200		{object}	object{page=view.Page}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/pages/{page} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	shown, err := logicdocs.Show(r.Runtime, c.Params("page"))
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"page": view.Page{
		Name: shown.Page,
		Text: shown.Text,
		HTML: shown.HTML,
	}})
}
