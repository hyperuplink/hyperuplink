package root

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/account"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/admin"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/attachment"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/categories"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/docs"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/newpost"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/search"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/session"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/user"
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
	r.Path = ""
	r.Env = route.NewEnv()

	r.Router.Route("/", func(base fiber.Router) {
		base.Get("", r.Index).Name("index")

		adminRoute, err := admin.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, adminRoute)

		sessionRoute, err := session.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionRoute)

		accountRoute, err := account.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, accountRoute)

		categoriesRoute, err := categories.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, categoriesRoute)

		docsRoute, err := docs.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, docsRoute)

		newpostRoute, err := newpost.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, newpostRoute)

		searchRoute, err := search.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, searchRoute)

		userRoute, err := user.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, userRoute)

		attachmentRoute, err := attachment.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, attachmentRoute)
	}, "root.")

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

func (r *Route) Show(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}

func (r *Route) Destroy(c fiber.Ctx) error {
	return c.SendStatus(fiber.StatusNotFound)
}
