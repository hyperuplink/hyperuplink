package logs

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
)

// @Summary		List the audit log
// @Description	The administrators' visits and settings changes, newest first.
// @Tags			admin
// @Produce		json
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Entries per page, at most 100"
// @Success		200			{object}	object{logs=[]view.Log,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		403			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/logs [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	page, perPage, ret, rerr := req.Page(req.System.GetTopicsPerPage())
	if ret {
		return rerr
	}

	logs, total, err := gh.Repositories(r.Runtime).Activity.VAllAdmin(common.QueryOptions{
		OrderBy: "a.created_at",
		Order:   common.Descending,
		Limit:   perPage,
		Page:    page,
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"logs":       view.Logs(logs),
		"pagination": view.NewPagination(page, perPage, total),
	})
}
