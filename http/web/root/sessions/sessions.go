package sessions

import (
	"github.com/gofiber/fiber/v3"
	"github.com/markbates/goth"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root/sessions/confirm"
	"github.com/mrusme/hyperuplink/runtime"

	"github.com/markbates/goth/providers/github"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"
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
	r.Path = route.For("Sessions").Pathname()
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
		base.Get("/"+route.For("SessionsSignin").Pathname(),
			r.SignInShow).Name("signin.show")
		base.Post("/"+route.For("SessionsSignin").Pathname(),
			r.SignInCreate).Name("signin.create")

		base.Get("/"+route.For("SessionsSignup").Pathname(),
			r.SignUpShow).Name("signup.show")
		base.Post("/"+route.For("SessionsSignup").Pathname(),
			r.SignUpCreate).Name("signup.create")

		base.Get("/"+route.For("SessionsSignout").Pathname(),
			r.SignOutShow).Name("signout.show")

		// base.Get("/tfa", r.TfaShow).Name("tfa.show")
		// base.Post("/tfa", r.TfaCreate).Name("tfa.create")

		// base.Get("/forgot", r.ForgotShow).Name("forgot.show")
		// base.Post("/forgot", r.ForgotCreate).Name("forgot.create")

		sessionsConfirmRoute, err := confirm.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionsConfirmRoute)
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
