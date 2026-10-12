package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/forums"
)

// @Summary		List the forums
// @Description	Every forum the viewer may read, ordered by category and position,
// @Description	each with the viewer's permissions in it, which is how a client
// @Description	finds out where it may post.
// @Tags			board
// @Produce		json
// @Success		200	{object}	object{forums=[]view.Forum}
// @Failure		401	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/forums [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	fums, err := logicforums.Index(req.Context(), r.Runtime, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"forums": req.Views().Forums(fums, req.Perms())})
}
