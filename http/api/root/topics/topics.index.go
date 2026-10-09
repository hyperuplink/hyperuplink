package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

// @Summary		List topics across the board
// @Description	One page of topics from every forum the viewer may read, pinned
// @Description	ones first, in the requested order: active (the default, newest
// @Description	activity first), new (newest topic first) or replies (most replies
// @Description	first).
// @Tags			board
// @Produce		json
// @Param			sort		query		string	false	"active, new or replies"
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Topics per page, at most 100"
// @Success		200			{object}	object{topics=[]view.Topic,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	page, perPage, ret, rerr := req.Page(req.System.GetTopicsPerPage())
	if ret {
		return rerr
	}
	sort, ret, rerr := req.Sort(string(topic.SortActive), topic.Sorts()...)
	if ret {
		return rerr
	}

	feed, err := topicsfeed.Index(r.Runtime, &topicsfeed.Input{
		Sort:     topic.Sort(sort),
		Page:     page,
		PerPage:  perPage,
		ViewerID: req.ViewerID(),
	}, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	tops := req.Views().Topics(*feed.Topics)
	view.AttachUnread(tops, feed.Unread)

	return req.Respond(fiber.Map{
		"topics":     tops,
		"pagination": view.NewPagination(page, perPage, feed.Total),
	})
}
