package attachment

import (
	"fmt"
	"path"
	"strings"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"xn--gckvb8fzb.com/hyperuplink/http/web/helpers"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

func (r *Route) Show(c fiber.Ctx) error {
	id, err := uuid.Parse(c.Params("attachment"))
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	att, err := r.Runtime.Repositories.Attachment.GetByUUID(id, common.QueryOptions{})
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	categoryID, found, err := r.Runtime.Repositories.Attachment.GetCategoryForAttachment(id)
	if err != nil || !found {
		return c.SendStatus(fiber.StatusNotFound)
	}

	role, memberOf := helpers.CurrentUserRoleAndGroups(r.Runtime, c)
	perms := helpers.ResolvePermissions(r.Runtime, role, memberOf)
	if !perms.CanReadID(categoryID) {
		return c.SendStatus(fiber.StatusNotFound)
	}

	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		r.Runtime.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}
	attachments := settingAttachments.JSONValue

	reader, err := r.Runtime.Storage.GetFile(
		attachments.StorageProviderID,
		path.Join(attachments.StoragePath, id.String()),
	)
	if err != nil {
		return c.SendStatus(fiber.StatusNotFound)
	}

	// Serving user-uploaded content like image/svg+xml or text/html inline
	// would allow scripts to execute and access the session cookie. Therefor
	// only serve images inline.
	disposition := "attachment"
	if attachments.InlineImageDisplay && inlineImageTypes[att.MimeType] {
		disposition = "inline"
	}

	c.Set("X-Content-Type-Options", "nosniff")
	if att.MimeType != "" {
		c.Set("Content-Type", att.MimeType)
	}
	c.Set("Content-Disposition", fmt.Sprintf(
		"%s; filename=\"%s\"",
		disposition,
		sanitizeFilename(att.Filename),
	))

	return c.SendStream(reader)
}

var inlineImageTypes = map[string]bool{
	"image/png":  true,
	"image/jpeg": true,
	"image/gif":  true,
	"image/webp": true,
}

func sanitizeFilename(name string) string {
	name = strings.Map(func(r rune) rune {
		if r < 0x20 || r == '"' || r == '\\' || r == '/' {
			return -1
		}
		return r
	}, name)
	if name == "" {
		return "attachment"
	}
	return name
}
