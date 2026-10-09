package api

import (
	"strings"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
)

func (srv *API) guests(c fiber.Ctx) error {
	sys := request.SystemFrom(c, srv.rt)

	if !request.PrincipalFrom(c).IsGuest() || sys.AllowAPIGuests || openForGuests(c) {
		return c.Next()
	}

	return problem.Respond(c, problem.From(errs.ErrAuthenticationRequired))
}

func openForGuests(c fiber.Ctx) bool {
	path := strings.TrimSuffix(c.Path(), "/")

	switch path {
	case "/" + route.For("APIV1").AsURL():
		return c.Method() == fiber.MethodGet
	case "/" + route.For("APIV1Session").AsURL():
		return c.Method() == fiber.MethodPost
	}

	return false
}
