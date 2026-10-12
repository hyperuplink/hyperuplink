package replies

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/helpers/notify"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/view"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	topicsfeed "xn--gckvb8fzb.com/hyperuplink/logic/root/topics"
)

// @Summary		Reply to a topic
// @Description	Adds a reply to the topic, or to one of its replies when parent_id
// @Description	names it. Attachments are uploaded beforehand and named in
// @Description	attachment_ids. The topic's author and earlier repliers are
// @Description	notified. The viewer needs write permission in the topic's
// @Description	category.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			topic	path		string					true	"The topic UUID or short id"
// @Param			request	body		replies.ReplyCreateBody	true	"The reply to store"
// @Success		201		{object}	object{reply=view.Reply}
// @Header			201		{string}	Location	"The path of the new reply"
// @Failure		401		{object}	problem.Problem
// @Failure		403		{object}	problem.Problem
// @Failure		404		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/topics/{topic}/replies [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	top, err := topicsfeed.Read(req.Context(), r.Runtime, c.Params("topic"), req.Perms())
	if ret, rerr := req.RespondOnReadError(err); ret {
		return rerr
	}

	in := new(ReplyCreateBody)
	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	created, err := logictopics.CreateReply(
		req.Context(),
		r.Runtime,
		req.User().ID,
		req.Perms(),
		&logictopics.CreateReplyInput{
			Text:    in.Text,
			TopicID: top.ID.String(),
			ReplyID: in.ParentID,
		},
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return logicattachments.ResolveOwned(req.Context(), r.Runtime, authorID, in.AttachmentIDs)
		},
	)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	notify.Reply(
		req.Context(),
		r.Runtime,
		created,
		req.User().Username,
		req.Ts("reply_notification_subject"),
	)

	vrep, err := gh.Repositories(r.Runtime).Reply.VGetByUUID(req.Context(), created.Reply.ID, common.QueryOptions{Limit: 1})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()
	m.PrefetchAttachments(vrep.AttachmentIDs)

	c.Set(fiber.HeaderLocation, view.APIPath("APIV1Reply", map[string]string{
		"reply": vrep.ID.String(),
	}))

	return req.RespondCreated(fiber.Map{
		"reply": m.Reply(vrep, top.CategorySlug, top.ForumSlug, top.Slug),
	})
}
