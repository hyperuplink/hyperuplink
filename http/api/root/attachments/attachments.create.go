package attachments

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
)

// @Summary		Upload attachments
// @Description	Stores one or more files sent as the multipart field "attachments"
// @Description	and returns them with their ids, which a topic or a reply then names
// @Description	in attachment_ids. Uploads are rejected when the board has
// @Description	attachments disabled.
// @Tags			board
// @Accept			mpfd
// @Produce		json
// @Param			attachments	formData	file	true	"One or more files"
// @Success		201			{object}	object{attachments=[]view.Attachment}
// @Failure		400			{object}	problem.Problem
// @Failure		401			{object}	problem.Problem
// @Failure		422			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/attachments [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	settings, err := logicattachments.Settings(req.Context(), r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}
	if !settings.EnableAttachments {
		return req.RespondError(errs.ErrAttachmentsDisabled)
	}

	form, ferr := c.MultipartForm()
	if ferr != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	files := form.File[logicattachments.FormField]
	if len(files) == 0 {
		return req.RespondProblem(problem.Validation(map[string][]string{
			logicattachments.FormField: {"required"},
		}))
	}

	ids, err := logicattachments.StoreAll(req.Context(), r.Runtime, req.User().ID, files)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	m := req.Views()
	m.PrefetchAttachments(ids)

	return req.RespondCreated(fiber.Map{"attachments": m.Attachments(ids)})
}
