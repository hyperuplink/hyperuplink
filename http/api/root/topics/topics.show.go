package topics

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
)

// @Summary		Show a topic
// @Description	One topic with its poll, addressed by UUID or short id, without its
// @Description	replies, which are a collection of their own. A signed-in viewer
// @Description	gets the read marker and the vote state as well. A topic the viewer
// @Description	may not read is reported as not found.
// @Tags			board
// @Produce		json
// @Param			topic	path		string	true	"The topic UUID or short id"
// @Success		200		{object}	object{topic=view.Topic}
// @Failure		401		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics/{topic} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	shown, err := topicsfeed.Show(req.Context(), r.Runtime, &topicsfeed.ShowInput{
		ID:       c.Params("topic"),
		ViewerID: req.ViewerID(),
	}, req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	m := req.Views()
	m.PrefetchAttachments(shown.Topic.AttachmentIDs)

	top := m.Topic(shown.Topic)
	top.Poll = view.PollFrom(shown.Poll, shown.Topic.EndedAt, !req.IsGuest())
	if !req.IsGuest() {
		top.Viewer = &view.TopicViewer{Unread: shown.Unread}
	}

	return req.Respond(fiber.Map{"topic": top})
}
