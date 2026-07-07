package search

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request/site"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vsearchresult"
	searchrepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/search"
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

	var query string = strings.TrimSpace(c.Query("q"))
	var author string = strings.TrimSpace(c.Query("author"))

	present := func(key string) bool { return c.Query(key) != "" }

	fTitle := present("f_title")
	fBody := present("f_body")
	fReplies := present("f_replies")
	fAttachments := present("f_attachments")
	if !fTitle && !fBody && !fReplies && !fAttachments {
		fTitle, fBody, fReplies = true, true, true
	}

	req.SetData("query", query)
	req.SetData("author", author)
	req.SetData("f_title", fTitle)
	req.SetData("f_body", fBody)
	req.SetData("f_replies", fReplies)
	req.SetData("f_attachments", fAttachments)

	if query == "" {
		return req.Respond()
	}

	var resultsPerPage int = req.System.GetTopicsPerPage()
	var postsPerPage int = req.System.GetPostsPerPage()

	var activePage int
	activePage, err = strconv.Atoi(c.Query("page", "1"))
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	var results *[]vsearchresult.VSearchResult
	var total int64
	results, total, err = r.Runtime.Repositories.Search.Query(
		query,
		searchrepo.Options{
			Title:                fTitle,
			Body:                 fBody,
			Replies:              fReplies,
			Attachments:          fAttachments,
			Author:               author,
			AllowedCategorySlugs: req.Perms().AllowedReadSlugs(),
		},
		resultsPerPage,
		activePage,
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	pages := helpers.GetNumberOfPages(total, resultsPerPage)
	req.Site.SetPager(site.NewPager(pages, resultsPerPage, activePage))

	req.SetData("results", results)
	req.SetData("total", int(total))
	req.SetData("posts_per_page", postsPerPage)

	extra := map[string]string{"q": query}
	if author != "" {
		extra["author"] = author
	}
	if fTitle {
		extra["f_title"] = "1"
	}
	if fBody {
		extra["f_body"] = "1"
	}
	if fReplies {
		extra["f_replies"] = "1"
	}
	if fAttachments {
		extra["f_attachments"] = "1"
	}
	req.SetData("pagination_extra", extra)

	return req.Respond()
}
