package profile

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logicprofile "xn--gckvb8fzb.com/hyperuplink/logic/root/account/profile"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
)

type Route struct {
	route.RouteController
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.For("AccountProfile").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("show")
		base.Put("", r.Update).Name("update")
		base.Post("/picture", r.PictureUpload).Name("picture")
	}, r.Path+".")

	return r, nil
}

func (r *Route) GetRuntime() *runtime.Runtime {
	return r.Runtime
}

func (r *Route) GetPath() string {
	return r.Path
}

func (r *Route) GetEnv() *route.Environment {
	return r.Env
}

// @Summary		Upload the account picture
// @Description	The picture is resized and converted by the board, so any format
// @Description	ImageMagick reads is accepted as long as it stays under the size
// @Description	limit the administrators have set.
// @Tags			account
// @Accept			mpfd
// @Produce		json
// @Param			profile_picture	formData	file	true	"The image to use as the account picture"
// @Success		200				{object}	request.StatusResponse
// @Failure		401				{object}	request.ErrorResponse
// @Failure		403				{object}	request.ErrorResponse
// @Failure		422				{object}	request.ErrorResponse
// @Security		BearerAuth
// @Security		APIKeyAuth
// @Router			/account/profile/picture [post]
func (r *Route) PictureUpload(c fiber.Ctx) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.UserRole,
		user.AdminRole,
	); ret {
		return rerr
	}

	var fh *multipart.FileHeader
	if form, ferr := c.MultipartForm(); ferr == nil {
		if files := form.File["profile_picture"]; len(files) > 0 {
			fh = files[0]
		}
	}
	if fh == nil {
		return req.RespondError(errs.ErrFormInvalid)
	}

	if err = logicprofile.StorePicture(r.Runtime, req.User, fh); err != nil {
		return req.RespondError(err)
	}

	err = r.Runtime.Repositories.User.Update(req.User)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
