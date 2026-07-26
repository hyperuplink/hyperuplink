package search

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	logicsearch "xn--gckvb8fzb.com/hyperuplink/logic/root/search"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	myRoute := route.For("Search")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.GuestRole, // TODO: Only if forum is public
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	present := func(key string) bool { return c.Query(key) != "" }

	var activePage int
	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var resultsPerPage int = req.System.GetTopicsPerPage()
	var postsPerPage int = req.System.GetPostsPerPage()

	in := &logicsearch.Input{
		Query:        strings.TrimSpace(c.Query("q")),
		Author:       strings.TrimSpace(c.Query("author")),
		FTitle:       present("f_title"),
		FBody:        present("f_body"),
		FReplies:     present("f_replies"),
		FAttachments: present("f_attachments"),
		Page:         activePage,
		PerPage:      resultsPerPage,
	}
	in.Normalize()

	req.SetData("query", in.Query)
	req.SetData("author", in.Author)
	req.SetData("f_title", in.FTitle)
	req.SetData("f_body", in.FBody)
	req.SetData("f_replies", in.FReplies)
	req.SetData("f_attachments", in.FAttachments)

	if in.Query == "" {
		return req.Respond()
	}

	var results *logicsearch.Results
	results, err = logicsearch.Query(r.Runtime, req.Perms(), in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	req.Site.SetPager(site.NewPager(results.Pages, resultsPerPage, activePage))

	req.SetData("results", results.Results)
	req.SetData("total", int(results.Total))
	req.SetData("posts_per_page", postsPerPage)

	extra := map[string]string{"q": in.Query}
	if in.Author != "" {
		extra["author"] = in.Author
	}
	if in.FTitle {
		extra["f_title"] = "1"
	}
	if in.FBody {
		extra["f_body"] = "1"
	}
	if in.FReplies {
		extra["f_replies"] = "1"
	}
	if in.FAttachments {
		extra["f_attachments"] = "1"
	}
	req.SetData("pagination_extra", extra)

	return req.Respond()
}
