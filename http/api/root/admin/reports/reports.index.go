package reports

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
)

// @Summary		List the reported posts
// @Tags			admin
// @Produce		json
// @Success		200	{object}	object{reports=[]view.Report}
// @Failure		401	{object}	problem.Problem
// @Failure		403	{object}	problem.Problem
// @Security		BearerAuth
// @Router			/admin/reports [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	reports, err := gh.Repositories(r.Runtime).PostEvent.VAllReports(common.QueryOptions{
		OrderBy: "pe.created_at",
		Order:   common.Descending,
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"reports": req.Views().Reports(*reports)})
}
