package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	logicthemes "xn--gckvb8fzb.com/hyperuplink/logic/root/admin/board/themes"
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
	r.Path = route.For("APIV1AdminBoardThemes").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Put("", r.Update).Name("update")
		base.Post("/"+route.For("APIV1AdminBoardThemesBanner").Pathname(), r.BannerUpload).Name("banner.upload")
		base.Delete("/"+route.For("APIV1AdminBoardThemesBanner").Pathname(), r.BannerRemove).Name("banner.remove")
		base.Post("/"+route.For("APIV1AdminBoardThemesFavicon").Pathname(), r.FaviconUpload).Name("favicon.upload")
		base.Delete("/"+route.For("APIV1AdminBoardThemesFavicon").Pathname(), r.FaviconRemove).Name("favicon.remove")
		base.Post("/"+route.For("APIV1AdminBoardThemesBackground").Pathname(), r.BackgroundUpload).Name("background.upload")
		base.Delete("/"+route.For("APIV1AdminBoardThemesBackground").Pathname(), r.BackgroundRemove).Name("background.remove")
	}, r.Path+".")

	return r, nil
}

const imageField = "image"

func (r *Route) imageUpload(c fiber.Ctx, kind logicthemes.ImageKind) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	form, ferr := c.MultipartForm()
	if ferr != nil {
		return req.RespondError(errs.ErrFormInvalid)
	}
	files := form.File[imageField]
	if len(files) == 0 {
		return req.RespondProblem(problem.Validation(map[string][]string{imageField: {"required"}}))
	}

	err = logicthemes.UploadImage(r.Runtime, kind, files[0])
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}

func (r *Route) imageRemove(c fiber.Ctx, kind logicthemes.ImageKind) (err error) {
	req := request.New(r, c)

	if ret, rerr := req.RequireAdmin(); ret {
		return rerr
	}

	err = logicthemes.RemoveImage(r.Runtime, kind)
	if ret, rerr := req.RespondOnError(err); ret {
		return rerr
	}

	return req.RespondNoContent()
}
