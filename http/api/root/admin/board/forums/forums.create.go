package forums

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/forums"
)

// @Summary		Create a forum
// @Tags			admin
// @Accept			json
// @Produce		json
// @Param			request	body		logicforums.CreateInput	true	"The forum to create"
// @Success		201		{object}	object{id=string}
// @Header			201		{string}	Location	"The path of the new forum"
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		409		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/forums [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	in := new(logicforums.CreateInput)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	id, err := logicforums.Create(req.Context(), r.Runtime, in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	c.Set(fiber.HeaderLocation, view.APIPath("APIV1AdminBoardForum", map[string]string{
		"forum": id.String(),
	}))

	return req.RespondCreated(fiber.Map{"id": id})
}
