package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/markbates/goth"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/site"
	"github.com/mrusme/hyperuplink/runtime"

	"github.com/markbates/goth/providers/github"
	goth_fiber "github.com/wakatara/goth_fiber"
)

type Route struct {
	route.Route
}

func New(
	rt *runtime.Runtime,
	router fiber.Router,
) (*Route, error) {
	r := new(Route)

	r.Runtime = rt
	r.Router = router
	r.Path = route.GetPathOf(route.SessionsRoute)
	r.Env = route.NewEnv()

	goth.UseProviders(
		github.New(
			"key",
			"secret",
			"callbackURL",
			"scope1", "scope2",
		),
	)

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("/signin", r.SignInShow).Name("signin.show")
		base.Post("/signin", r.SignInCreate).Name("signin.create")
		base.Get("/:provider", goth_fiber.BeginAuthHandler).Name("provider.show")
		base.Get("/:provider/callback", func(ctx fiber.Ctx) error {
			user, err := goth_fiber.CompleteUserAuth(ctx)
			if err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			return ctx.SendString(user.Email)
		}).Name("provider.callback.show")
		base.Get("/signout", func(ctx fiber.Ctx) error {
			if err := goth_fiber.Logout(ctx); err != nil {
				return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
			}
			return ctx.SendString("logged out")
		}).Name("signout.show")
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

func (r *Route) SignInShow(c fiber.Ctx) error {
	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}

func (r *Route) SignInCreate(c fiber.Ctx) error {
	return c.Render("views/session/signin", fiber.Map{
		"Site": site.New(r, c),
	}, "views/layouts/base")
}

func (r *Route) Index(c fiber.Ctx) error {
	return c.SendString("I'm a INDEX request!")
}

func (r *Route) Show(c fiber.Ctx) error {
	return c.SendString("I'm a SHOW request!")
}

func (r *Route) Create(c fiber.Ctx) error {
	return c.SendString("I'm a CREATE request!")
}

func (r *Route) Update(c fiber.Ctx) error {
	return c.SendString("I'm a UPDATE request!")
}

func (r *Route) Destroy(c fiber.Ctx) error {
	return c.SendString("I'm a DESTROY request!")
}
