package api

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/validation"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories/common"
)

type API struct {
	rt        *runtime.Runtime
	app       *fiber.App
	r         route.IRouteController
	validator *validator.Validate
}

func New(
	rt *runtime.Runtime,
) (srv *API, err error) {
	srv = new(API)
	srv.rt = rt

	srv.validator = validation.New()

	srv.app = fiber.New(fiber.Config{
		StrictRouting:      false,
		CaseSensitive:      false,
		BodyLimit:          srv.rt.Config.APIBodyLimit(),
		Concurrency:        srv.rt.Config.APIConcurrency(),
		ProxyHeader:        srv.rt.Config.APIProxyHeader(),
		EnableIPValidation: srv.rt.Config.APIEnableIPValidation(),
		TrustProxy:         srv.rt.Config.APITrustProxy(),
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback: srv.rt.Config.APITrustLoopback(),
			Proxies:  srv.rt.Config.APITrustProxies(),
		},
		ReduceMemoryUsage: srv.rt.Config.APIReduceMemoryUsage(),
		ServerHeader:      srv.rt.Config.APIServerHeader(),
		AppName:           "hyperuplink-api",
		StructValidator:   validation.NewStructValidator(srv.validator),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if ferr, ok := err.(*fiber.Error); ok {
				code = ferr.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	return srv, nil
}

func (srv *API) loadMiddlewares() error {
	srv.app.Use(slogfiber.NewWithConfig(srv.rt.Logger, slogfiber.Config{
		DefaultLevel:       srv.rt.LoggerLevel,
		WithRequestID:      true,
		WithRequestBody:    srv.rt.IsDevelopmentMode(),
		WithRequestHeader:  true,
		WithResponseBody:   false,
		WithResponseHeader: true,
		WithSpanID:         true,
		WithTraceID:        srv.rt.IsDevelopmentMode(),
	}))
	srv.app.Use(recover.New())
	srv.app.Use(requestid.New())
	srv.app.Use(cors.New())

	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.LivenessEndpoint),
		healthcheck.New())
	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.ReadinessEndpoint),
		healthcheck.New())
	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.StartupEndpoint),
		healthcheck.New())

	srv.app.Use(srv.authenticate)

	return nil
}

func (srv *API) authenticate(c fiber.Ctx) error {
	secret := extractSecret(c)
	if secret == "" || !apikey.IsSecret(secret) {
		return unauthorized(c)
	}

	key, err := srv.rt.Repositories.APIKey.GetBySecretHash(
		apikey.HashSecret(secret),
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return unauthorized(c)
	}

	usr, err := srv.rt.Repositories.User.GetByID(
		key.UserID.String(),
		common.QueryOptions{
			WithBanned:  false,
			WithSpammed: false,
			WithDeleted: false,
			Limit:       1,
		},
	)
	if err != nil {
		return unauthorized(c)
	}

	if terr := srv.rt.Repositories.APIKey.TouchLastUsed(key.ID); terr != nil {
		srv.rt.Warn("error", terr)
	}

	c.Locals(request.UserLocal, usr)

	return c.Next()
}

func extractSecret(c fiber.Ctx) string {
	if header := c.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}

	return c.Get("X-API-Key")
}

func unauthorized(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"error": errs.ErrAPIKeyInvalid.Error(),
	})
}

func (srv *API) loadRoutes() (err error) {
	if srv.r, err = root.New(srv.rt, srv.app); err != nil {
		return err
	}

	return nil
}

func (srv *API) Startup() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.loadMiddlewares(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	if err = srv.loadRoutes(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *API) Run() error {
	listenAddr := fmt.Sprintf(
		"%s:%d",
		srv.rt.Config.APIBindIP(),
		srv.rt.Config.APIPort(),
	)
	if err := srv.app.Listen(listenAddr, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil &&
		err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (srv *API) Shutdown() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.app.ShutdownWithTimeout(time.Second * 5); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}
