package picture

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicaccount "xn--gckvb8fzb.com/hyperuplink/logic/root/account"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
)

const formField = "picture"

// @Summary		Upload the profile picture
// @Description	Replaces the picture with the image sent as the multipart field
// @Description	"picture", converted to the board's picture format. Rejected when
// @Description	the board has pictures disabled.
// @Tags			account
// @Accept			mpfd
// @Produce		json
// @Param			picture	formData	file	true	"The image"
// @Success		200		{object}	object{user=view.Account}
// @Failure		400		{object}	problem.Problem
// @Failure		401		{object}	problem.Problem
// @Failure		422		{object}	problem.Problem
// @Security		BearerAuth
// @Router			/account/picture [post]
func (r *Route) Create(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireUser(); ret {
		return rerr
	}

	enabled, err := logicprofile.PicturesEnabled(req.Context(), r.Runtime)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}
	if !enabled {
		return req.RespondError(errs.ErrPicturesDisabled)
	}

	var fh *multipart.FileHeader
	form, ferr := c.MultipartForm()
	if ferr != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}
	if files := form.File[formField]; len(files) > 0 {
		fh = files[0]
	}
	if fh == nil {
		return req.RespondProblem(problem.Validation(map[string][]string{formField: {"required"}}))
	}

	usr := req.User()
	if err = logicprofile.StorePicture(req.Context(), r.Runtime, usr, fh); err != nil {
		return req.RespondError(err)
	}

	err = gh.Repositories(r.Runtime).User.Update(req.Context(), usr)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	fresh, err := gh.Repositories(r.Runtime).User.GetByUUID(req.Context(), usr.ID, common.QueryOptions{Limit: 1})
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	acc, err := logicaccount.Show(req.Context(), r.Runtime, fresh)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.Respond(fiber.Map{"user": req.Views().Account(fresh, acc.Groups, acc.Profile)})
}
