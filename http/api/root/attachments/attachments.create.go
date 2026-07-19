package attachments

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
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

	authorID, ok := req.UserUUID()
	if !ok {
		return req.RespondError(errs.ErrForbidden)
	}

	form, ferr := c.MultipartForm()
	if ferr != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	files := form.File[logicattachments.FormField]
	if len(files) == 0 {
		return req.RespondError(errs.ErrFormInvalid)
	}

	ids, err := logicattachments.StoreAll(r.Runtime, authorID, files)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondCreated(fiber.Map{
		"ids": ids,
	})
}
