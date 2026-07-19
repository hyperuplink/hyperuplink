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
