package attachment

import (
	"fmt"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	logicperms "xn--gckvb8fzb.com/hyperuplink/logic/helpers/perms"
	logicattachment "xn--gckvb8fzb.com/hyperuplink/logic/root/attachment"
)

func (r *Route) Show(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("attachment"))
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	role, memberOf := helpers.CurrentUserRoleAndGroups(r.Runtime, c)
	perms := logicperms.Resolve(r.Runtime, role, memberOf)

	file, err := logicattachment.Open(r.Runtime, id, perms)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
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
