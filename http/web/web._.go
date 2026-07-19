package web

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"os"
	runt "runtime"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/sprig/v3"
	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/extractors"
	"github.com/gofiber/fiber/v3/middleware/cache"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/csrf"
	"github.com/gofiber/fiber/v3/middleware/healthcheck"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/gofiber/fiber/v3/middleware/requestid"
	goth_fiber "github.com/shareed2k/goth_fiber/v2"

	"github.com/gofiber/fiber/v3/middleware/session"
	"github.com/gofiber/fiber/v3/middleware/static"

	"github.com/gofiber/storage/redis/v3"
	html "github.com/gofiber/template/html/v3"
	slogfiber "github.com/samber/slog-fiber"
	"xn--gckvb8fzb.com/hyperuplink/errs"
	"xn--gckvb8fzb.com/hyperuplink/http/route"
	"xn--gckvb8fzb.com/hyperuplink/http/validation"
	"xn--gckvb8fzb.com/hyperuplink/http/web/root"
	"xn--gckvb8fzb.com/hyperuplink/models/setting"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	settingRepo "xn--gckvb8fzb.com/hyperuplink/services/repositories/setting"
)

type Web struct {
	rt        *runtime.Runtime
	app       *fiber.App
	r         route.IRouteController
	engine    *html.Engine
	validator *validator.Validate
	watcher   *fsnotify.Watcher
	hash      string
}

func New(
	rt *runtime.Runtime,
) (srv *Web, err error) {
	srv = new(Web)
	srv.rt = rt

	if srv.engine, err = srv.getViewsEngine(); err != nil {
		return nil, err
	}

	srv.validator = validation.New()

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
		Views:             srv.engine,
		StructValidator:   validation.NewStructValidator(srv.validator),
	})

	return srv, nil
}

func (srv *Web) loadMiddlewares() error {
	srv.app.Use(slogfiber.NewWithConfig(srv.rt.Logger, slogfiber.Config{
		DefaultLevel:       srv.rt.LoggerLevel,
		WithRequestID:      true,
		WithRequestBody:    srv.rt.IsDevelopmentMode(),
		WithRequestHeader:  true,
		WithResponseBody:   false, // srv.rt.IsDevelopmentMode(),
		WithResponseHeader: true,
		WithSpanID:         true,
		WithTraceID:        srv.rt.IsDevelopmentMode(),
		Filters:            []slogfiber.Filter{slogfiber.IgnorePathPrefix("/static/")},
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

	if srv.rt.IsDevelopmentMode() == false {
		// TODO: Adjust config
		srv.app.Use(cache.New(cache.ConfigDefault))
	}

	storage, err := srv.getSessionStorage()
	if err != nil {
		return err
	}

	// In development mode cookies are served over plain HTTP, so the "Secure"
	// (HTTPS-only) flag must be disabled or the session/CSRF cookies would never
	// be sent back by the browser (or curl)
	secureCookies := srv.rt.IsDevelopmentMode() == false

	sessConfig := session.Config{
		Storage:         storage,
		CookieSecure:    secureCookies,  // HTTPS only (disabled in development)
		CookieHTTPOnly:  true,           // Prevent XSS
		CookieSameSite:  "Lax",          // CSRF protection
		IdleTimeout:     6 * time.Hour,  // Session timeout
		AbsoluteTimeout: 24 * time.Hour, // Maximum session life
		Extractor:       extractors.FromCookie("__hyperuplink_session"),
	}

	sess, sessStore := session.NewWithStore(sessConfig)

	srv.app.Use(sess)

	goth_fiber.SessionManager = goth_fiber.NewSessionManager(sessStore)

	srv.app.Use(csrf.New(csrf.Config{
		CookieName:            "__hyperuplink_csrf",
		CookieSecure:          secureCookies, // HTTPS only (disabled in development)
		CookieHTTPOnly:        true,          // Needs to be 'false' to allow JS to access tokens
		CookieSameSite:        "Lax",
		CookieSessionOnly:     true,
		Extractor:             extractors.FromForm("_csrf"),
		Session:               sessStore,
		DisableValueRedaction: true,
	}))

	return nil
}

func (srv *Web) getSessionStorage() (storage fiber.Storage, err error) {
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

	storage = redis.New(redisConfig)

	return storage, nil
}

func (srv *Web) getViewsEngine() (*html.Engine, error) {
	var engine *html.Engine

	if srv.rt.IsDevelopmentMode() {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		engine = html.New(cwd, ".html")

		srv.rt.Debug("cwd", cwd)
		srv.watcher, err = srv.getWatcher(cwd+"/views", func(f string) {
			srv.rt.Debug("change", f)
			err := srv.app.ReloadViews()
			if err != nil {
				srv.rt.Error("error", err)
			}
		})
	} else {
		engine = html.NewFileSystem(http.FS(srv.rt.Embeds["views"]), ".html")
	}
	// https://masterminds.github.io/sprig/
	engine.AddFuncMap(sprig.FuncMap())
	engine.AddFunc("safeHTML",
		func(s string) template.HTML {
			return template.HTML(s)
		})
	engine.AddFunc("safeURL",
		func(s string) template.URL {
			return template.URL(s)
		})

	return engine, nil
}

func (srv *Web) getWatcher(path string, callback func(f string)) (*fsnotify.Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case ev, ok := <-watcher.Events:
				if !ok {
					return
				}
				if ev.Has(fsnotify.Write) {
					callback(ev.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				srv.rt.Error("error", err)
			}
		}
	}()

	if err = watcher.Add(path); err != nil {
		watcher.Close()
		return nil, err
	}

	return watcher, nil
}

func (srv *Web) loadRoutes() (err error) {
	var stic fiber.Handler

	if srv.rt.IsDevelopmentMode() {
		cwd, err := os.Getwd()
		if err != nil {
			return err
		}
		stic = static.New(cwd+"/static", static.Config{
			Browse:        true,
			CacheDuration: -1,
		})
	} else {
		var sub fs.FS
		if sub, err = fs.Sub(srv.rt.Embeds["static"], "static"); err != nil {
			return err
		}
		stic = static.New("", static.Config{
			FS:            sub,
			Browse:        false,
			CacheDuration: 10 * time.Second, // TODO: Make configurable
		})
	}
	srv.app.Get("/static*", stic)

	if err = srv.loadStorageRoutes(); err != nil {
		return err
	}

	if srv.r, err = root.New(srv.rt, srv.app); err != nil {
		return err
	}

	return nil
}

func (srv *Web) loadStorageRoutes() (err error) {
	storages, err := srv.rt.Config.Storages()
	if err != nil {
		return err
	}
	for _, storageCfg := range storages {
		if strings.ToLower(storageCfg.Type) != "local" {
			continue
		}
		if storageCfg.Local.Path == "" || storageCfg.Local.PublicURI == "" {
			continue
		}
		if route.CollidesWithRoute(storageCfg.Local.PublicURI) {
			return fmt.Errorf(
				"local storage %q PublicURI %q collides with a route",
				storageCfg.ID,
				storageCfg.Local.PublicURI,
			)
		}
		providerID := storageCfg.ID
		staticHandler := static.New(
			storageCfg.Local.Path,
			static.Config{
				Browse:        false,
				CacheDuration: 10 * time.Second, // TODO: Make configurable
			},
		)
		srv.app.Get(storageCfg.Local.PublicURI+"*", func(c fiber.Ctx) error {
			if srv.isGatedAttachmentPath(providerID, c.Params("*")) {
				return c.SendStatus(fiber.StatusNotFound)
			}
			return staticHandler(c)
		})
	}

	return nil
}

// isGatedAttachmentPath reports whether a direct static request targets an
// attachment file. Attachment files are only served through the permission
// checked /attachment/:id route, so direct access to their storage path must be
// blocked to prevent bypassing the read permission check by guessing URLs.
func (srv *Web) isGatedAttachmentPath(providerID, rel string) bool {
	settingAttachments, err := settingRepo.GetByID[setting.Attachments](
		srv.rt.Repositories.Setting,
		"attachments",
	)
	if err != nil {
		return false
	}
	att := settingAttachments.JSONValue

	if att.StorageProviderID != providerID || att.StoragePath == "" {
		return false
	}

	rel = strings.TrimPrefix(rel, "/")
	prefix := strings.Trim(att.StoragePath, "/") + "/"

	return strings.HasPrefix(rel, prefix)
}

func (srv *Web) Startup() (err error) {
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

func (srv *Web) Shutdown() (err error) {
	srv.rt.Debug("status", "exec")

	if err = srv.app.ShutdownWithTimeout(time.Second * 5); err != nil {
		srv.rt.Error("status", "error", "error", err)
		return err
	}

	srv.rt.Info("status", "ok")

	return nil
}
