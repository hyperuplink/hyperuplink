package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/forums"
)

// @Summary		Show a forum
// @Description	One forum, addressed by UUID or slug. A forum the viewer may not
// @Description	read is reported as not found.
// @Tags			board
// @Produce		json
// @Param			forum	path		string	true	"The forum UUID or slug"
// @Success		200		{object}	object{forum=view.Forum}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/forums/{forum} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	fum, err := logicforums.Show(r.Runtime, c.Params("forum"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"forum": req.Views().Forum(fum, req.Perms())})
}
