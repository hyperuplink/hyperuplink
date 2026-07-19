package topics

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	logictopics "xn--gckvb8fzb.com/hyperuplink/logic/root/categories/forums/topics"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

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

	r.notifyReply(c, req, created)

	return req.RespondCreated(fiber.Map{
		"id":       created.Reply.ID,
		"short_id": created.Reply.ShortID,
	})
}
