package session

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/session/confirm"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root/session/provider"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	repoSetting "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
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
	r.Path = route.For("Session").Pathname()
	r.Env = route.NewEnv()

	baseURL := ""
	if settingSystem, serr := repoSetting.GetByID[setting.System](
		rt.Repositories.Setting,
		"system",
	); serr == nil {
		baseURL = settingSystem.JSONValue.BaseURL
	} else {
		rt.Warn("auth_provider", "system setting unavailable", "error", serr)
	}

	if err := RegisterAuthProviders(rt, baseURL); err != nil {
		return nil, err
	}

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Get("/",
			r.Index).Name("index")

		base.Get("/"+route.For("SessionSignin").Pathname(),
			r.SignInShow).Name("signin.show")
		base.Post("/"+route.For("SessionSignin").Pathname(),
			r.SignInCreate).Name("signin.create")

		base.Get("/"+route.For("SessionSignup").Pathname(),
			r.SignUpShow).Name("signup.show")
		base.Post("/"+route.For("SessionSignup").Pathname(),
			r.SignUpCreate).Name("signup.create")

		base.Get("/"+route.For("SessionSignout").Pathname(),
			r.SignOutShow).Name("signout.show")

		base.Get("/"+route.For("SessionTwofactor").Pathname(),
			r.TfaShow).Name("twofactor.show")
		base.Post("/"+route.For("SessionTwofactor").Pathname(),
			r.TfaCreate).Name("twofactor.create")

		// base.Get("/forgot", r.ForgotShow).Name("forgot.show")
		// base.Post("/forgot", r.ForgotCreate).Name("forgot.create")

		sessionConfirmRoute, err := confirm.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionConfirmRoute)
		// Warning: Do not add routes below this point, as :provider will have
		// preference over them. Add any fixed route before this line.
		sessionProviderRoute, err := provider.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, sessionProviderRoute)
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
