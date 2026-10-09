package api

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/etag"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/http/validation"
	"xn--gckvb8fzb.com/glides/runtime"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/problem"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root"
	"xn--gckvb8fzb.com/hyperuplink/http/helpers/redisstorage"
)

const healthPrefix = "/_internal/health"

type API struct {
	rt        *runtime.Runtime
	cfg       gh.APIConfig
	app       *fiber.App
	r         route.IRouteController
	validator *validator.Validate
	storage   fiber.Storage
}

// @title			Hyperuplink API
// @version		{{.Version}}
// @description	The JSON interface to a Hyperuplink board. It serves the
// @description	categories, forums, topics and replies the board shows in the
// @description	browser, the account of the caller and the administration, under
// @description	the prefix /api/v1 on a port of its own.
// @description
// @description	A request without a key is a guest and reads what the board's
// @description	default permission grants to guests, as long as the administrator
// @description	allows guest access to the API. A request with a key acts as the
// @description	user the key belongs to. Keys are obtained by signing in through
// @description	POST /session with username and password, or on the board under
// @description	Account -> API, and they are sent as a bearer token in the
// @description	Authorization header. Errors are RFC 9457 problem details with a
// @description	stable machine code in the "code" member.
//
// @contact.name	Hyperuplink
// @contact.url	https://hyperup.link
//
// @license.name	SEGV
// @license.url	https://xn--gckvb8fzb.com/segv/
//
// @basePath	/api/v1
// @accept		json
// @produce	json
//
// @tag.name		board
// @tag.description	The board: its categories, forums, topics, replies and users.
// @tag.name		session
// @tag.description	Signing in, the viewer and the key of the request.
// @tag.name		account
// @tag.description	The signed-in user's own profile, settings, keys and credentials.
// @tag.name		admin
// @tag.description	Board administration, restricted to administrators.
//
// @securityDefinitions.apikey	BearerAuth
// @in				header
// @name			Authorization
// @description		The API key, as "Bearer hup_...".
func New(
	rt *runtime.Runtime,
	cfg gh.APIConfig,
) (srv *API, err error) {
	srv = new(API)
	srv.rt = rt
	srv.cfg = cfg

	srv.validator = validation.New()

	srv.app = fiber.New(fiber.Config{
		StrictRouting:      false,
		CaseSensitive:      false,
		BodyLimit:          cfg.BodyLimit,
		Concurrency:        cfg.Concurrency,
		ProxyHeader:        cfg.ProxyHeader,
		EnableIPValidation: cfg.EnableIPValidation,
		TrustProxy:         cfg.TrustProxy,
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback: cfg.TrustLoopback,
			Proxies:  cfg.TrustProxies,
		},
		ReduceMemoryUsage: cfg.ReduceMemoryUsage,
		ServerHeader:      cfg.ServerHeader,
		AppName:           "hyperuplink-api",
		StructValidator:   validation.NewStructValidator(srv.validator),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			return problem.Respond(c, problem.FromFiber(err))
		},
	})

	return srv, nil
}

func (srv *API) App() *fiber.App {
	return srv.app
}

func (srv *API) loadMiddlewares() (err error) {
	srv.app.Use(slogfiber.NewWithConfig(srv.rt.Logger(), slogfiber.Config{
		DefaultLevel:       srv.rt.GetLogLevel(),
		WithRequestID:      true,
		WithRequestBody:    false,
		WithRequestHeader:  true,
		WithResponseBody:   false,
		WithResponseHeader: true,
		WithSpanID:         true,
		WithTraceID:        srv.rt.IsDevelopmentMode(),
	}))
	srv.app.Use(recover.New())
	srv.app.Use(requestid.New())

	if len(srv.cfg.CORSOrigins) > 0 {
		srv.app.Use(cors.New(cors.Config{
			AllowOrigins: srv.cfg.CORSOrigins,
			AllowMethods: []string{
				fiber.MethodGet,
				fiber.MethodPost,
				fiber.MethodPut,
				fiber.MethodPatch,
				fiber.MethodDelete,
				fiber.MethodOptions,
			},
			AllowHeaders: []string{
				fiber.HeaderAuthorization,
				fiber.HeaderContentType,
			},
			ExposeHeaders: []string{
				fiber.HeaderETag,
				fiber.HeaderRetryAfter,
				fiber.HeaderWWWAuthenticate,
				"X-Request-Id",
			},
		}))
	}

	srv.app.Use(cacheControl)

	srv.app.Get(healthPrefix+healthcheck.LivenessEndpoint, healthcheck.New())
	srv.app.Get(healthPrefix+healthcheck.ReadinessEndpoint, healthcheck.New())
	srv.app.Get(healthPrefix+healthcheck.StartupEndpoint, healthcheck.New())

	srv.loadOpenAPI()

	srv.app.Use(srv.principal)
	srv.app.Use(srv.guests)

	if srv.storage, err = redisstorage.New(srv.rt.Config()); err != nil {
		return err
	}
	srv.app.Use(srv.limiter())
	srv.app.Use("/"+route.For("APIV1Session").AsURL(), srv.authLimiter())

	srv.app.Use(etag.New(etag.Config{
		Weak: true,
		Next: func(c fiber.Ctx) bool {
			return c.Method() != fiber.MethodGet
		},
	}))

	return nil
}

func cacheControl(c fiber.Ctx) error {
	c.Set(fiber.HeaderCacheControl, "private, no-cache")

	return c.Next()
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
	listenAddr := fmt.Sprintf("%s:%d", srv.cfg.BindIP, srv.cfg.Port)
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
