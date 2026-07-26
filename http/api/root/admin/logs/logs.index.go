package logs

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary	List the audit log
// @Tags		admin
// @Produce	json
// @Param		page	query		integer	false	"The page to return, counting from one"
// @Success	200		{object}	object{logs=[]vactivity.VActivity,total=integer,pages=integer}
// @Failure	401		{object}	request.ErrorResponse
// @Failure	403		{object}	request.ErrorResponse
// @Failure	422		{object}	request.ErrorResponse
// @Security	BearerAuth
// @Security	APIKeyAuth
// @Router		/admin/logs [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var activePage int
	if activePage, err = strconv.Atoi(c.Query("page", "1")); err != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	perPage := req.System.GetTopicsPerPage()

	logs, total, err := gh.Repositories(r.Runtime).Activity.VAllAdmin(
		common.QueryOptions{
			OrderBy: "a.created_at",
			Order:   common.Descending,
			Limit:   perPage,
			Page:    activePage,
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(fiber.Map{
		"logs":  logs,
		"total": total,
		"pages": paging.Pages(total, perPage),
	})
}
