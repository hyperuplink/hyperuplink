package themes

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
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
	r.Path = route.For("AdminBoardThemes").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("", r.Index).Name("index")
		base.Post("/update", r.Update).Name("update")
		base.Post("/banner", r.BannerUpload).Name("banner")
		base.Post("/banner/remove", r.BannerRemove).Name("banner.remove")
		base.Post("/favicon", r.FaviconUpload).Name("favicon")
		base.Post("/favicon/remove", r.FaviconRemove).Name("favicon.remove")
		base.Post("/background", r.BackgroundUpload).Name("background")
		base.Post("/background/remove", r.BackgroundRemove).Name("background.remove")
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
