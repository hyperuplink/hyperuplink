package newpost

import (
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	logicnewpost "xn--gckvb8fzb.com/hyperuplink/logic/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/models/vtopic"
)

// @Summary		Start a topic
// @Description	Creates a topic in the given forum, optionally with a poll and with
// @Description	attachments that were uploaded beforehand, and answers with the slugs
// @Description	the new topic can be read under.
// @Tags			board
// @Accept			json
// @Produce		json
// @Param			request	body		newpost.NewCreateBody	true	"The topic to create"
// @Success		201		{object}	object{id=string,slug=string,category_slug=string,forum_slug=string}
// @Failure		401		{object}	request.ErrorResponse
// @Failure		403		{object}	request.ErrorResponse
// @Failure		422		{object}	request.ValidationErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/new [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	in := new(NewCreateBody)

	if ret, rerr := req.BindJSON(in); ret {
		return rerr
	}

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	location, lerr := time.LoadLocation(req.User.Timezone)
	if lerr != nil {
		location = time.UTC
	}

	var vtop *vtopic.VTopic
	vtop, err = logicnewpost.Create(
		r.Runtime,
		authorID,
		req.Perms(),
		location,
		&in.CreateInput,
		func(authorID uuid.UUID) ([]uuid.UUID, error) {
			return logicattachments.ResolveOwned(r.Runtime, authorID, in.AttachmentIDs)
		},
	)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"id":            vtop.ID,
		"slug":          vtop.Slug,
		"category_slug": vtop.CategorySlug,
		"forum_slug":    vtop.ForumSlug,
	})
}
