package docs

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	logicdocs "xn--gckvb8fzb.com/hyperuplink/logic/root/docs"
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
	r.Path = route.For("Docs").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("/"+logicdocs.PageAbout, r.page(logicdocs.PageAbout)).Name("about")
		base.Get("/"+logicdocs.PageContact, r.page(logicdocs.PageContact)).Name("contact")
		base.Get("/"+logicdocs.PagePrivacy, r.page(logicdocs.PagePrivacy)).Name("privacy")
		base.Get("/"+logicdocs.PageTerms, r.page(logicdocs.PageTerms)).Name("terms")
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

func (r *Route) page(page string) fiber.Handler {
	return func(c fiber.Ctx) (err error) {
		req := request.New(r, c)

		if ret, rerr := req.AccessControl(
			user.UserRole,
			user.AdminRole,
		); ret {
			return rerr
		}

		var view *logicdocs.Page
		view, err = logicdocs.Show(r.Runtime, page)
		if ret, rerr := req.RespondOnError(err); ret == true {
			return rerr
		}

		return req.Respond(view)
	}
}
