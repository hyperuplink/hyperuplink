package logs

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

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

	logs, total, err := r.Runtime.Repositories.Activity.VAllAdmin(
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
