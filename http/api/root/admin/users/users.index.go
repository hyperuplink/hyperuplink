package users

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
)

// @Summary		List the users
// @Description	Every user including the banned and the deleted ones, newest first.
// @Tags			admin
// @Produce		json
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Users per page, at most 100"
// @Success		200			{object}	object{users=[]view.AdminUser,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		403			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/users [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	page, perPage, ret, rerr := req.Page(req.System.GetTopicsPerPage())
	if ret {
		return rerr
	}

	all := common.QueryOptions{WithBanned: true, WithDeleted: true}

	total, err := gh.Repositories(r.Runtime).User.CountAll(req.Context(), all)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	all.OrderBy = "created_at"
	all.Order = common.Descending
	all.Limit = perPage
	all.Page = page

	users, err := gh.Repositories(r.Runtime).User.All(req.Context(), all)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"users":      req.Views().AdminUsers(users),
		"pagination": view.NewPagination(page, perPage, total),
	})
}
