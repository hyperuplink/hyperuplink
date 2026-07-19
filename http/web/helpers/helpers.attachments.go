package helpers

import (
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

const AttachmentsFormField = logicattachments.FormField

func ProcessAttachments(
	rt *runtime.Runtime,
	c fiber.Ctx,
	authorID uuid.UUID,
) (ids []uuid.UUID, err error) {
	form, ferr := c.MultipartForm()
	if ferr != nil {
		return nil, nil
	}

	return logicattachments.StoreAll(rt, authorID, form.File[AttachmentsFormField])
}
