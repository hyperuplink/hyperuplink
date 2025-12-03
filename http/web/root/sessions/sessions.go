package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/markbates/goth"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/runtime"

	"github.com/markbates/goth/providers/github"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
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

		base.Get("/signup", r.SignUpShow).Name("signup.show")
		base.Post("/signup", r.SignUpCreate).Name("signup.create")

		base.Get("/signout", r.SignOutShow).Name("signout.show")

		base.Get("/confirm", r.ConfirmShow).Name("confirm.show")
		base.Post("/confirm", r.ConfirmCreate).Name("confirm.create")

		// base.Get("/tfa", r.TfaShow).Name("tfa.show")
		// base.Post("/tfa", r.TfaCreate).Name("tfa.create")

		// base.Get("/forgot", r.ForgotShow).Name("forgot.show")
		// base.Post("/forgot", r.ForgotCreate).Name("forgot.create")

		// Warning: Do not add routes below this point, as :provider will have
		// preference over them. Add any fixed route before this line.
		base.Get("/:provider", goth_fiber.BeginAuthHandler).Name("provider.show")
		base.Get("/:provider/callback", r.ProviderCallbackShow).Name("provider.callback.show")
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
