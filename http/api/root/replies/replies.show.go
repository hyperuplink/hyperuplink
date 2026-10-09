package replies

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	"xn--gckvb8fzb.com/hyperuplink/logic/helpers/paging"
	logicreplies "xn--gckvb8fzb.com/hyperuplink/logic/root/replies"
)

// @Summary		Show a reply
// @Description	One reply, addressed by UUID or short id, with the topic it belongs
// @Description	to, its position among the topic's replies and the page of the
// @Description	reply collection it is on at the board's posts-per-page setting.
// @Tags			board
// @Produce		json
// @Param			reply	path		string	true	"The reply UUID or short id"
// @Success		200		{object}	object{reply=view.Reply,topic=view.Ref,position=integer,page=integer}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/replies/{reply} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	shown, err := logicreplies.Show(r.Runtime, c.Params("reply"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	m := req.Views()
	m.PrefetchAttachments(shown.Reply.AttachmentIDs)

	return req.Respond(fiber.Map{
		"reply": m.Reply(shown.Reply, shown.Topic.CategorySlug, shown.Topic.ForumSlug, shown.Topic.Slug),
		"topic": view.Ref{
			ID:   shown.Topic.ID.String(),
			Name: shown.Topic.Name,
			Slug: shown.Topic.Slug,
		},
		"position": shown.Position,
		"page":     paging.Pages(shown.Position, req.System.GetPostsPerPage()),
	})
}
