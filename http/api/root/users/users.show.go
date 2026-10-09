package users

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicuser "xn--gckvb8fzb.com/hyperuplink/logic/root/user"
)

// @Summary		Show a user
// @Description	The public profile of a user with the groups they belong to and
// @Description	their ten newest topics and replies the viewer may read.
// @Tags			board
// @Produce		json
// @Param			user	path		string	true	"The username"
// @Success		200		{object}	object{user=view.User,topics=[]view.Topic,replies=[]view.Reply}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/users/{user} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	shown, err := logicuser.Show(r.Runtime, c.Params("user"), req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()

	ids := make([][]uuid.UUID, 0, len(shown.Replies))
	for i := range shown.Replies {
		ids = append(ids, shown.Replies[i].AttachmentIDs)
	}
	m.PrefetchAttachments(ids...)

	replies := make([]view.Reply, 0, len(shown.Replies))
	for i := range shown.Replies {
		rep := &shown.Replies[i]
		replies = append(replies, m.Reply(&rep.VReply, rep.CategorySlug, rep.ForumSlug, rep.TopicSlug))
	}

	return req.Respond(fiber.Map{
		"user":    m.User(shown.User, shown.Membership),
		"topics":  m.Topics(shown.Topics),
		"replies": replies,
	})
}
