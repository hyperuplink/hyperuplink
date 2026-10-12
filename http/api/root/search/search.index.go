package search

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicsearch "xn--gckvb8fzb.com/hyperuplink/logic/root/search"
)

const (
	inTitle   = "title"
	inBody    = "body"
	inReplies = "replies"
)

// @Summary		Search the board
// @Description	Full-text search over what the viewer may read. The "in" parameter
// @Description	narrows the match to topic titles, topic bodies and replies in any
// @Description	combination; without it all three are searched.
// @Tags			board
// @Produce		json
// @Param			q			query		string	true	"The terms to search for"
// @Param			author		query		string	false	"Only match posts by this username"
// @Param			in			query		string	false	"A comma-separated subset of title, body and replies"
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Results per page, at most 100"
// @Success		200			{object}	object{results=[]view.SearchResult,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/search [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	fields := map[string][]string{}

	query := strings.TrimSpace(c.Query("q"))
	if query == "" {
		fields["q"] = []string{"required"}
	} else if len(query) > 256 {
		fields["q"] = []string{"max"}
	}

	author := strings.TrimSpace(c.Query("author"))
	if len(author) > 32 {
		fields["author"] = []string{"max"}
	}

	in := &logicsearch.Input{Query: query, Author: author}
	if raw := strings.TrimSpace(c.Query("in")); raw != "" {
		for _, scope := range strings.Split(raw, ",") {
			switch strings.TrimSpace(scope) {
			case inTitle:
				in.FTitle = true
			case inBody:
				in.FBody = true
			case inReplies:
				in.FReplies = true
			default:
				fields["in"] = []string{"oneof"}
			}
		}
	}

	if len(fields) > 0 {
		return req.RespondProblem(problem.Validation(fields))
	}

	page, perPage, ret, rerr := req.Page(req.System.GetTopicsPerPage())
	if ret {
		return rerr
	}
	in.Page = page
	in.PerPage = perPage
	in.Normalize()

	results, err := logicsearch.Query(r.Runtime, req.Perms(), in)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"results":    req.Views().SearchResults(results.Results, req.System.GetPostsPerPage()),
		"pagination": view.NewPagination(page, perPage, results.Total),
	})
}
