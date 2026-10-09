package api

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/limiter"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
)

const (
	limitWindow    = time.Minute
	limitKeyPrefix = "ratelimit:"
)

func (srv *API) limitFor(c fiber.Ctx) int {
	if request.PrincipalFrom(c).IsGuest() {
		return srv.cfg.RateLimitGuest
	}

	return srv.cfg.RateLimitUser
}

func limitKey(c fiber.Ctx) string {
	if p := request.PrincipalFrom(c); p.Key != nil {
		return limitKeyPrefix + "key:" + p.Key.ID.String()
	}

	return limitKeyPrefix + "ip:" + c.IP()
}

func limitReached(c fiber.Ctx) error {
	return problem.Respond(c, problem.From(errs.ErrRateLimited))
}

func authKey(c fiber.Ctx) string {
	var body struct {
		Username string `json:"username"`
	}
	username := "-"
	if json.Unmarshal(c.Body(), &body) == nil && body.Username != "" {
		username = strings.ToLower(strings.TrimSpace(body.Username))
	}

	return limitKeyPrefix + "auth:" + c.IP() + ":" + username
}

func (srv *API) authLimiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Storage:    srv.storage,
		Expiration: limitWindow,
		Next: func(c fiber.Ctx) bool {
			return c.Method() != fiber.MethodPost || srv.cfg.RateLimitAuth <= 0
		},
		MaxFunc: func(c fiber.Ctx) int {
			return srv.cfg.RateLimitAuth
		},
		KeyGenerator: authKey,
		LimitReached: limitReached,
	})
}

func (srv *API) limiter() fiber.Handler {
	return limiter.New(limiter.Config{
		Storage:    srv.storage,
		Expiration: limitWindow,
		Next: func(c fiber.Ctx) bool {
			return srv.limitFor(c) <= 0
		},
		MaxFunc:      srv.limitFor,
		KeyGenerator: limitKey,
		LimitReached: limitReached,
	})
}
