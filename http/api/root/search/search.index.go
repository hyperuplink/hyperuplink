package search

import (
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicsearch "xn--gckvb8fzb.com/hyperuplink/logic/root/search"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	present := func(key string) bool { return c.Query(key) != "" }

	var activePage int
	if activePage, err = strconv.Atoi(c.Query("page", "1")); err != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	in := &logicsearch.Input{
		Query:        strings.TrimSpace(c.Query("q")),
		Author:       strings.TrimSpace(c.Query("author")),
		FTitle:       present("f_title"),
		FBody:        present("f_body"),
		FReplies:     present("f_replies"),
		FAttachments: present("f_attachments"),
		Page:         activePage,
		PerPage:      req.System.GetTopicsPerPage(),
	}
	in.Normalize()

	if in.Query == "" {
		return req.RespondError(errs.ErrValidation)
	}

	var results *logicsearch.Results
	results, err = logicsearch.Query(r.Runtime, req.Perms(), in)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.Respond(results)
}
