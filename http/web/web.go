package web

import (
	"fmt"
	"net/http"
	runt "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/storage/redis/v3"
	"github.com/mrusme/hyperuplink/errs"
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
		// ErrorHandler: func(c fiber.Ctx, err error) error {
		// 	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		// 		"errors":  []string{err.Error()},
		// 		"status":  0,
		// 		"request": requestid.FromContext(c),
		// 	})
		// },
	})

	return srv, nil
}

func (srv *Web) LoadMiddlewares() error {
	srv.app.Use(slogfiber.NewWithConfig(srv.rt.Logger, slogfiber.Config{
		DefaultLevel:       srv.rt.LoggerLevel,
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  true,
		WithResponseBody:   false,
		WithResponseHeader: true,
		WithSpanID:         true,
		WithTraceID:        true,
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

	// ---------------------------------------------------------------------------

	storage, err := srv.LoadSessionStorage()
	if err != nil {
		return err
	}

	sessConfig := session.Config{
		Storage:         storage,
		CookieSecure:    true,             // HTTPS only
		CookieHTTPOnly:  true,             // Prevent XSS
		CookieSameSite:  "Lax",            // CSRF protection
		IdleTimeout:     30 * time.Minute, // Session timeout
		AbsoluteTimeout: 24 * time.Hour,   // Maximum session life
		Extractor:       extractors.FromCookie("__Host-session_id"),
	}

	srv.app.Use(session.New(sessConfig))

	srv.app.Use(csrf.New(csrf.Config{
		CookieName:        "__Host-csrf_",
		CookieSecure:      true,
		CookieHTTPOnly:    true, // Needs to be 'false' to allow JS to access tokens
		CookieSameSite:    "Lax",
		CookieSessionOnly: true,
		Extractor:         extractors.FromHeader("X-Csrf-Token"),
		// Storage:           storage,
		Session: session.NewStore(sessConfig),
		// DisableValueRedaction: true,
	}))
	return nil
}

func (srv *Web) LoadSessionStorage() (fiber.Storage, error) {
	var err error

	redisConfig := redis.Config{
		MasterName: srv.rt.Config.RedisMasterName(),
		Username:   srv.rt.Config.RedisUsername(),
		Password:   srv.rt.Config.RedisPassword(),
		Database:   srv.rt.Config.RedisDatabase(),
		Reset:      srv.rt.Config.RedisReset(),
	}

	addrs := srv.rt.Config.RedisAddresses()
	addrsl := len(addrs)
	if addrsl == 0 {
		return nil, errs.ErrRedisAddrsEmpty
	} else if addrsl == 1 {
		hostPort := strings.Split(addrs[0], ":")
		if len(hostPort) != 2 {
			return nil, errs.ErrRedisAddrsMalformed
		}
		redisConfig.Host = hostPort[0]
		if redisConfig.Port, err = strconv.Atoi(hostPort[1]); err != nil {
			return nil, errs.ErrRedisAddrsMalformed
		}
	} else if addrsl > 1 {
		redisConfig.Addrs = addrs
	}

	poolSize := srv.rt.Config.RedisPoolsize()
	if poolSize <= 0 {
		poolSize = 10 * runt.GOMAXPROCS(0)
	}
	redisConfig.PoolSize = poolSize

	storage := redis.New(redisConfig)

	return storage, nil
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
