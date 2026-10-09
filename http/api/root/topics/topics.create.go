package topics

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
)

// @Summary		Start a topic
// @Description	Creates a topic in the forum named by forum_id, as a regular topic
// @Description	or as a poll with poll_options and an optional poll_ends_at in the
// @Description	author's timezone. Attachments are uploaded beforehand and named in
// @Description	attachment_ids. The viewer needs write permission in the forum's
// @Description	category.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			request	body		topics.TopicCreateBody	true	"The topic to create"
// @Success		201		{object}	object{topic=view.Topic}
// @Header			201		{string}	Location	"The path of the new topic"
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	in := new(TopicCreateBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	location, lerr := time.LoadLocation(req.User().Timezone)
	if lerr != nil {
		location = time.UTC
	}

	vtop, err := logicnewpost.Create(
		r.Runtime,
		req.User().ID,
		req.Perms(),
		location,
		&in.CreateInput,
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return logicattachments.ResolveOwned(r.Runtime, authorID, in.AttachmentIDs)
		},
	)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()
	m.PrefetchAttachments(vtop.AttachmentIDs)

	top := m.Topic(vtop)
	if vtop.IsPoll() {
		var poll *logictopics.Poll
		poll, err = logictopics.PollView(r.Runtime, &logictopics.PollViewInput{
			Topic:    &vtop.Topic,
			ViewerID: req.ViewerID(),
			CanWrite: true,
		})
		if ret, rerr := req.RespondOnError(err); ret {
			return rerr
		}
		top.Poll = view.PollFrom(poll, vtop.EndedAt, true)
	}

	c.Set(fiber.HeaderLocation, view.APIPath("APIV1Topic", map[string]string{
		"topic": vtop.ID.String(),
	}))

	return req.RespondCreated(fiber.Map{"topic": top})
}
