package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicforums "xn--gckvb8fzb.com/hyperuplink/logic/root/forums"
	"xn--gckvb8fzb.com/hyperuplink/models/topic"
)

// @Summary		List a forum's topics
// @Description	One page of the forum's topics, pinned ones first, in the requested
// @Description	order: active (the default, newest activity first), new (newest
// @Description	topic first) or replies (most replies first).
// @Tags			board
// @Produce		json
// @Param			forum		path		string	true	"The forum UUID or slug"
// @Param			sort		query		string	false	"active, new or replies"
// @Param			page		query		integer	false	"The page, counting from one"
// @Param			per_page	query		integer	false	"Topics per page, at most 100"
// @Success		200			{object}	object{forum=view.Forum,topics=[]view.Topic,pagination=view.Pagination}
// @Failure		401			{object}	problem.Problem
// @Failure		404			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/forums/{forum}/topics [get]
func (r *Route) Index(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	fum, err := logicforums.Show(r.Runtime, c.Params("forum"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	page, perPage, ret, rerr := req.Page(req.System.GetTopicsPerPage())
	if ret {
		return rerr
	}
	sort, ret, rerr := req.Sort(string(topic.SortActive), topic.Sorts()...)
	if ret {
		return rerr
	}

	listed, err := logicforums.ListTopics(r.Runtime, &logicforums.TopicsInput{
		Forum:    fum,
		Sort:     topic.Sort(sort),
		Page:     page,
		PerPage:  perPage,
		ViewerID: req.ViewerID(),
	})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()
	tops := m.Topics(*listed.Topics)
	view.AttachUnread(tops, listed.Unread)

	return req.Respond(fiber.Map{
		"forum":      m.Forum(fum, req.Perms()),
		"topics":     tops,
		"pagination": view.NewPagination(page, perPage, listed.Total),
	})
}
