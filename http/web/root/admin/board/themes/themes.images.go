package themes

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
)

func (r *Route) imageUpload(
	c fiber.Ctx,
	kind logicthemes.ImageKind,
	formField string,
) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	var fh *multipart.FileHeader
	if form, ferr := c.MultipartForm(); ferr == nil {
		if files := form.File[formField]; len(files) > 0 {
			fh = files[0]
		}
	}

	if err = logicthemes.UploadImage(r.Runtime, kind, fh); err != nil {
		req.Flash.SetError(err)
	}

	return req.RedirectToRoute(myRoute)
}

func (r *Route) imageRemove(
	c fiber.Ctx,
	kind logicthemes.ImageKind,
) (err error) {
	myRoute := route.For("AdminBoardThemes")
	req := request.New(r, c, myRoute,
		[]string{"base"}, myRoute.AsURL()+"/index",
		myRoute.AsTitle())

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	if err = logicthemes.RemoveImage(r.Runtime, kind); err != nil {
		req.Flash.SetError(err)
	}

	return req.RedirectToRoute(myRoute)
}
