package categories

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
)

// @Summary		List the categories
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{categories=[]view.Category}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/board/categories [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	cats, err := gh.Repositories(r.Runtime).Category.All(common.QueryOptions{
		OrderBy: "position",
		Order:   common.Ascending,
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()
	views := make([]view.Category, 0, len(*cats))
	for i := range *cats {
		views = append(views, m.Category(&(*cats)[i], nil, req.Perms()))
	}

	return req.Respond(fiber.Map{"categories": views})
}
