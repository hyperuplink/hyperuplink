package web

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/mrusme/hyperuplink/runtime"
	slogfiber "github.com/samber/slog-fiber"
)

type Web struct {
	rt  *runtime.Runtime
	app *fiber.App
}

func New(
	rt *runtime.Runtime,
) (*Web, error) {
	srv := new(Web)

	srv.rt = rt
	srv.app = fiber.New(fiber.Config{
		StrictRouting:      false,
		CaseSensitive:      false,
		BodyLimit:          srv.rt.Config.ServerBodyLimit(),
		Concurrency:        srv.rt.Config.ServerConcurrency(),
		ProxyHeader:        srv.rt.Config.ServerProxyHeader(),
		EnableIPValidation: srv.rt.Config.ServerEnableIPValidation(),
		TrustProxy:         srv.rt.Config.ServerTrustProxy(),
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback: srv.rt.Config.ServerTrustLoopback(),
			Proxies:  srv.rt.Config.ServerTrustProxies(),
		},
		ReduceMemoryUsage: srv.rt.Config.ServerReduceMemoryUsage(),
		ServerHeader:      srv.rt.Config.ServerServerHeader(),
		AppName:           "hyperuplink",
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"errors":  []string{err.Error()},
				"status":  0,
				"request": requestid.FromContext(c),
			})
		},
	})

	return srv, nil
}

func (srv *Web) LoadMiddlewares() error {
	srv.app.Use(slogfiber.New(srv.rt.Logger))
	srv.app.Use(recover.New())
	srv.app.Use(requestid.New())
	srv.app.Use(cors.New())

	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.LivenessEndpoint),
		healthcheck.New())
	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.ReadinessEndpoint),
		healthcheck.New())
	srv.app.Get(fmt.Sprintf("/_internal/health%s", healthcheck.StartupEndpoint),
		healthcheck.New())

	return nil
}

func (srv *Web) Startup() error {
	var err error

	srv.rt.Debug("status", "exec")

	if err = srv.LoadMiddlewares(); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}

func (srv *Web) Run() error {
	listenAddr := fmt.Sprintf(
		"%s:%d",
		srv.rt.Config.ServerBindIP(),
		srv.rt.Config.ServerPort(),
	)
	if err := srv.app.Listen(listenAddr, fiber.ListenConfig{
		DisableStartupMessage: true,
	}); err != nil &&
		err != http.ErrServerClosed {
		return err
	}

	return nil
}

func (srv *Web) Shutdown() error {
	var err error

	srv.rt.Debug("status", "exec")

	if err = srv.app.ShutdownWithTimeout(time.Second * 5); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}
