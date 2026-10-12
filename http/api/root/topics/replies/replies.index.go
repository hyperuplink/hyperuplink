package replies

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicreplies "xn--gckvb8fzb.com/hyperuplink/logic/root/replies"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
)

// @Summary		List a topic's replies
// @Description	One page of the topic's replies, oldest first. A reply that answers
// @Description	another reply names it in parent_id.
// @Tags			board
// @Produce		json
// @Param			topic		path		string	true	"The topic UUID or short id"
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Replies per page, at most 100"
// @Success		200			{object}	object{replies=[]view.Reply,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		404			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics/{topic}/replies [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	top, err := topicsfeed.Read(r.Runtime, c.Params("topic"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	page, perPage, ret, rerr := req.Page(req.System.GetPostsPerPage())
	if ret {
		return rerr
	}

	listed, err := logicreplies.Index(r.Runtime, &logicreplies.IndexInput{
		Topic:   top,
		Page:    page,
		PerPage: perPage,
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{
		"replies":    req.Views().Replies(listed.Replies, top.CategorySlug, top.ForumSlug, top.Slug),
		"pagination": view.NewPagination(page, perPage, listed.Total),
	})
}
