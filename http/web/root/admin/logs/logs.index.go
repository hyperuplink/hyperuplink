package logs

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vactivity"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("AdminLogs")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var activePage int
	var perPage int = req.System.GetTopicsPerPage()

	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var logs *[]vactivity.VActivity
	var total int64
	logs, total, err = gh.Repositories(r.Runtime).Activity.VAllAdmin(
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

	pages := paging.Pages(total, perPage)
	req.Site.SetPager(site.NewPager(pages, perPage, activePage))

	req.SetData("logs", logs)

	return req.Respond()
}
