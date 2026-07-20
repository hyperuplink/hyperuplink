package attachments

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachment "xn--gckvb8fzb.com/hyperuplink/logic/root/attachment"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Download an attachment
// @Description	The file is streamed back under the name it was uploaded with, and
// @Description	callers who cannot read the category the attachment was posted in
// @Description	are answered as though it did not exist.
// @Tags			board
// @Produce		octet-stream
// @Param			attachment	path		string	true	"The attachment identifier"
// @Success		200			{file}		binary
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		404			{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/attachments/{attachment} [get]
func (r *Route) Show(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	id, err := uuid.Parse(c.Params("attachment"))
	if err != nil {
		return req.RespondError(errs.ErrNoRows)
	}

	file, err := logicattachment.Open(r.Runtime, id, req.Perms())
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	c.Set("X-Content-Type-Options", "nosniff")
	if file.Attachment.MimeType != "" {
		c.Set("Content-Type", file.Attachment.MimeType)
	}
	c.Set("Content-Disposition", fmt.Sprintf(
		"%s; filename=\"%s\"",
		file.Disposition,
		logicattachment.SanitizeFilename(file.Attachment.Filename),
	))

	return c.SendStream(file.Reader)
}
