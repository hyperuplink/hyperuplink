package api

import (
	"errors"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

const (
	bearerScheme     = "Bearer"
	lastUsedInterval = time.Minute
)

func (srv *API) principal(c fiber.Ctx) error {
	header := c.Get(fiber.HeaderAuthorization)
	if header == "" {
		return c.Next()
	}

	scheme, token, found := strings.Cut(header, " ")
	token = strings.TrimSpace(token)
	if !found || !strings.EqualFold(scheme, bearerScheme) || !apikey.IsSecret(token) {
		return problem.Respond(c, problem.From(errs.ErrAPIKeyInvalid))
	}

	key, err := gh.Repositories(srv.rt).APIKey.GetBySecretHash(
		apikey.HashSecret(token),
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return srv.rejectKey(c, err)
	}

	usr, err := gh.Repositories(srv.rt).User.GetByID(
		key.UserID.String(),
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return srv.rejectKey(c, err)
	}

	srv.touchLastUsed(key)

	c.Locals(request.PrincipalLocal, request.Principal{User: usr, Key: key})

	return c.Next()
}

func (srv *API) rejectKey(c fiber.Ctx, err error) error {
	if errors.Is(err, errs.ErrNoRows) {
		return problem.Respond(c, problem.From(errs.ErrAPIKeyInvalid))
	}

	srv.rt.Error("error", err)

	return problem.Respond(c, problem.From(err))
}

func (srv *API) touchLastUsed(key *apikey.APIKey) {
	if key.LastUsedAt.Valid && time.Since(key.LastUsedAt.Time) < lastUsedInterval {
		return
	}

	if err := gh.Repositories(srv.rt).APIKey.TouchLastUsed(key.ID); err != nil {
		srv.rt.Warn("error", err)
	}
}
