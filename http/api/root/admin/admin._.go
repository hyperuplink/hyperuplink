package admin

import (
	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/auth"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/board"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/comms"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/general"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/health"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/logs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/permissions"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/reports"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root/admin/users"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
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
	r.Path = route.For("APIV1Admin").Pathname()
	r.Env = route.NewEnv()

	r.Router.Route("/"+r.Path, func(base fiber.Router) {
		base.Use(r.RecordVisit)

		authRoute, err := auth.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, authRoute)

		boardRoute, err := board.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, boardRoute)

		commsRoute, err := comms.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, commsRoute)

		generalRoute, err := general.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, generalRoute)

		healthRoute, err := health.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, healthRoute)

		logsRoute, err := logs.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, logsRoute)

		permissionsRoute, err := permissions.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, permissionsRoute)

		reportsRoute, err := reports.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, reportsRoute)

		usersRoute, err := users.New(r.Runtime, base)
		r.Runtime.NilOrDie(err)
		r.Routes = append(r.Routes, usersRoute)
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

func (r *Route) RecordVisit(c fiber.Ctx) (err error) {
	if c.Method() != fiber.MethodGet {
		return c.Next()
	}

	principal := request.PrincipalFrom(c)
	if principal.User == nil || principal.User.Role != user.AdminRole {
		return c.Next()
	}

	logicactivity.RecordAdminVisit(r.Runtime, principal.User.ID, c.Path())

	return c.Next()
}
