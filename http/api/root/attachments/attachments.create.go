package attachments

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicattachments "xn--gckvb8fzb.com/hyperuplink/logic/helpers/attachments"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

// @Summary		Upload attachments
// @Description	Files are uploaded on their own and come back as identifiers, which
// @Description	are then passed as attachment_ids when the post that carries them is
// @Description	created.
// @Tags			board
// @Accept			mpfd
// @Produce		json
// @Param			attachments	formData	file	true	"One or more files to store"
// @Success		201			{object}	object{ids=[]string}
// @Failure		401			{object}	request.ErrorResponse
// @Failure		403			{object}	request.ErrorResponse
// @Failure		422			{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/attachments [post]
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
