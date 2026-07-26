package themes

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
	"xn--gckvb8fzb.com/hyperuplink/models/user"
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
	r.Path = route.For("AdminBoardThemes").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Put("", r.Update).Name("update")
		base.Post("/banner", r.BannerUpload).Name("banner")
		base.Delete("/banner", r.BannerRemove).Name("banner.remove")
		base.Post("/favicon", r.FaviconUpload).Name("favicon")
		base.Delete("/favicon", r.FaviconRemove).Name("favicon.remove")
		base.Post("/background", r.BackgroundUpload).Name("background")
		base.Delete("/background", r.BackgroundRemove).Name("background.remove")
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

func (r *Route) imageUpload(
	c fiber.Ctx,
	kind logicthemes.ImageKind,
	formField string,
) (err error) {
	req := request.New(r, c)

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

	err = logicthemes.UploadImage(r.Runtime, kind, fh)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}

func (r *Route) imageRemove(
	c fiber.Ctx,
	kind logicthemes.ImageKind,
) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.AccessControl(
		user.AdminRole,
	); ret {
		return rerr
	}

	err = logicthemes.RemoveImage(r.Runtime, kind)
	if ret, rerr := req.RespondOnError(err); ret == true {
		return rerr
	}

	return req.RespondOK()
}
