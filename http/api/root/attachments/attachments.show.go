package attachments

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachment "xn--gckvb8fzb.com/hyperuplink/logic/root/attachment"
)

// @Summary		Download an attachment
// @Description	The file is streamed back under the name it was uploaded with. An
// @Description	attachment of a post the viewer may not read is reported as not
// @Description	found.
// @Tags			board
// @Produce		octet-stream
// @Param			attachment	path		string	true	"The attachment UUID"
// @Success		200			{file}		binary
// @Failure		401			{object}	problem.Problem
// @Failure		404			{object}	problem.Problem
// @Security		BearerAuth
// @Router			/attachments/{attachment} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	id, err := uuid.Parse(c.Params("attachment"))
	if err != nil {
		return req.RespondError(errs.ErrNoRows)
	}

	file, err := logicattachment.Open(req.Context(), r.Runtime, id, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	c.Set("X-Content-Type-Options", "nosniff")
	if file.Attachment.MimeType != "" {
		c.Set(fiber.HeaderContentType, file.Attachment.MimeType)
	}
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf(
		"%s; filename=\"%s\"",
		file.Disposition,
		logicattachment.SanitizeFilename(file.Attachment.Filename),
	))

	return c.SendStream(file.Reader)
}
