package reports

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	List the reported posts
// @Tags		admin
// @Produce	json
// @Success	200	{object}	object{reports=[]vpostevent.VPostEvent}
// @Failure	401	{object}	request.ErrorResponse
// @Failure	403	{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/reports [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	reports, err := gh.Repositories(r.Runtime).PostEvent.VAllReports(common.QueryOptions{
		OrderBy: "pe.created_at",
		Order:   common.Descending,
	})
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"reports": reports,
	})
}
