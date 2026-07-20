package topics

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/helpers/notify"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Reply to a topic
// @Description	The reply is addressed by the slugs in the path, and everyone
// @Description	subscribed to the topic is notified once it has been stored.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			categories	path		string					true	"The category slug"
// @Param			forums		path		string					true	"The forum slug"
// @Param			topics		path		string					true	"The topic slug"
// @Param			request		body		topics.TopicCreateBody	true	"The reply to store"
// @Success		201			{object}	object{id=string,short_id=string}
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Failure		422			{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/_{categories}/{forums}/{topics} [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(TopicCreateBody)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	var created *logictopics.CreatedReply
	created, err = logictopics.CreateReply(
		r.Runtime,
		authorID,
		req.Perms(),
		&in.CreateReplyInput,
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return logicattachments.ResolveOwned(r.Runtime, authorID, in.AttachmentIDs)
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	notify.Reply(
		r.Runtime,
		created,
		req.User.Username,
		req.Ts("reply_notification_subject"),
	)

	return req.RespondCreated(fiber.Map{
		"id":       created.Reply.ID,
		"short_id": created.Reply.ShortID,
	})
}
