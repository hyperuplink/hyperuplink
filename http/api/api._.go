package api

import (
	"fmt"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	slogfiber "github.com/samber/slog-fiber"
	"github.com/swaggo/swag"
	"xn--gckvb8fzb.com/glides/http/route"
	"xn--gckvb8fzb.com/glides/http/validation"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/services/repositories/common"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api/request"
	"xn--gckvb8fzb.com/hyperuplink/http/api/root"
	"xn--gckvb8fzb.com/hyperuplink/models/apikey"
)

const (
	internalPrefix  = "/_internal/"
	swaggerSpecPath = "docs/swagger.json"
)

type API struct {
	rt        *runtime.Runtime
	app       *fiber.App
	r         route.IRouteController
	validator *validator.Validate
}

// @title		HyperUplink API
// @version		{{.Version}}
// @description	The JSON interface to a HyperUplink board, exposing the same
// @description	categories, forums, topics and administrative surfaces the web
// @description	interface serves, with every response shaped by the permissions
// @description	of the user the API key belongs to.
// @description
// @description	Authentication happens with an API key that a signed-in user
// @description	mints under Account -> API on the board itself, and the secret
// @description	it hands out once is sent either as a bearer token in the
// @description	Authorization header or verbatim in the X-API-Key header, so a
// @description	request that carries neither is answered with 401 before it
// @description	ever reaches a controller.
//
// @contact.name	HyperUplink
// @contact.url	https://hyperup.link
//
// @license.name	SEGV
// @license.url	https://xn--gckvb8fzb.com/segv/
//
// @basePath	/
// @accept		json
// @produce		json
//
// @tag.name		board
// @tag.description	The board itself, its categories, forums and topics.
// @tag.name		account
// @tag.description	The signed-in user's own profile, settings and credentials.
// @tag.name		admin
// @tag.description	Board administration, restricted to administrators.
//
// @securityDefinitions.apikey	BearerAuth
// @in				header
// @name			Authorization
// @description		The API key secret prefixed with "Bearer ".
//
// @securityDefinitions.apikey	APIKeyAuth
// @in				header
// @name			X-API-Key
// @description		The API key secret on its own.
func New(
	rt *runtime.Runtime,
) (srv *API, err error) {
	srv = new(API)
	srv.rt = rt

	srv.validator = validation.New()

	srv.app = fiber.New(fiber.Config{
		StrictRouting:      false,
		CaseSensitive:      false,
		BodyLimit:          srv.rt.Config().APIBodyLimit(),
		Concurrency:        srv.rt.Config().APIConcurrency(),
		ProxyHeader:        srv.rt.Config().APIProxyHeader(),
		EnableIPValidation: srv.rt.Config().APIEnableIPValidation(),
		TrustProxy:         srv.rt.Config().APITrustProxy(),
		TrustProxyConfig: fiber.TrustProxyConfig{
			Loopback: srv.rt.Config().APITrustLoopback(),
			Proxies:  srv.rt.Config().APITrustProxies(),
		},
		ReduceMemoryUsage: srv.rt.Config().APIReduceMemoryUsage(),
		ServerHeader:      srv.rt.Config().APIServerHeader(),
		AppName:           "hyperuplink-api",
		StructValidator:   validation.NewStructValidator(srv.validator),
		ErrorHandler: func(c fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if ferr, ok := err.(*fiber.Error); ok {
				code = ferr.Code
			}
			return c.Status(code).JSON(request.ErrorResponse{Error: err.Error()})
		},
	})

	return srv, nil
}

func (srv *API) loadMiddlewares() error {
	srv.app.Use(slogfiber.NewWithConfig(srv.rt.Logger(), slogfiber.Config{
		DefaultLevel:       srv.rt.GetLogLevel(),
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

	srv.loadSwagger()

	srv.app.Use(srv.authenticate)

	return nil
}

func (srv *API) loadSwagger() {
	embedded, ok := srv.rt.GetEmbedOk("docs")
	if !ok {
		srv.rt.Warn("swagger", "no embedded documentation")
		return
	}

	doc, err := fs.ReadFile(embedded, swaggerSpecPath)
	if err != nil {
		srv.rt.Warn("swagger", "error", "error", err)
		return
	}

	version, _, _, _ := srv.rt.GetBuild()
	swag.Register(swag.Name, &swag.Spec{
		Version:         version,
		SwaggerTemplate: string(doc),
	})

	srv.app.Get("/_internal/swagger/*", swaggo.HandlerDefault)
}

func (srv *API) authenticate(c fiber.Ctx) error {
	if isInternal(c) {
		return c.Next()
	}

	secret := extractSecret(c)
	if secret == "" || !apikey.IsSecret(secret) {
		return unauthorized(c)
	}

	key, err := gh.Repositories(srv.rt).APIKey.GetBySecretHash(
		apikey.HashSecret(secret),
		common.QueryOptions{Limit: 1},
	)
	if err != nil {
		return unauthorized(c)
	}

	usr, err := gh.Repositories(srv.rt).User.GetByID(
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

	if terr := gh.Repositories(srv.rt).APIKey.TouchLastUsed(key.ID); terr != nil {
		srv.rt.Warn("error", terr)
	}

	c.Locals(request.UserLocal, usr)

	return c.Next()
}

func isInternal(c fiber.Ctx) bool {
	return strings.HasPrefix(strings.ToLower(c.Path()), internalPrefix)
}

func extractSecret(c fiber.Ctx) string {
	if header := c.Get("Authorization"); strings.HasPrefix(header, "Bearer ") {
		return strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	}

	return c.Get("X-API-Key")
}

func unauthorized(c fiber.Ctx) error {
	return c.Status(fiber.StatusUnauthorized).JSON(request.ErrorResponse{
		Error: errs.ErrAPIKeyInvalid.Error(),
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
		srv.rt.Config().APIBindIP(),
		srv.rt.Config().APIPort(),
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
