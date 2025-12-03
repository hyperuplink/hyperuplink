package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	runt "runtime"
	"strconv"
	"strings"
	"time"

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
	"github.com/mrusme/hyperuplink/errs"
	"github.com/mrusme/hyperuplink/http/route"
	"github.com/mrusme/hyperuplink/http/web/root"
	"github.com/mrusme/hyperuplink/runtime"
	slogfiber "github.com/samber/slog-fiber"
)

type Web struct {
	rt      *runtime.Runtime
	app     *fiber.App
	r       route.IRoute
	engine  *html.Engine
	watcher *fsnotify.Watcher
	hash    string
}

type structValidator struct {
	validate *validator.Validate
}

func (v *structValidator) Validate(out any) error {
	return v.validate.Struct(out)
}

func New(
	rt *runtime.Runtime,
) (srv *Web, err error) {
	srv = new(Web)
	srv.rt = rt

	if srv.engine, err = srv.getViewsEngine(); err != nil {
		return nil, err
	}

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
		StructValidator:   &structValidator{validate: validator.New()},
	})

	return srv, nil
}

func (srv *Web) loadMiddlewares() error {
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

	if srv.rt.IsDevelopmentMode() == false {
		// TODO: Adjust config
		srv.app.Use(cache.New(cache.ConfigDefault))
	}

	storage, err := srv.getSessionStorage()
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
		Extractor:       extractors.FromCookie("__hyperuplink_session"),
	}

	sess, sessStore := session.NewWithStore(sessConfig)

	srv.app.Use(sess)

	goth_fiber.SessionManager = goth_fiber.NewSessionManager(sessStore)

	srv.app.Use(csrf.New(csrf.Config{
		CookieName:            "__hyperuplink_csrf",
		CookieSecure:          true,
		CookieHTTPOnly:        true, // Needs to be 'false' to allow JS to access tokens
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
	engine.AddFunc(
		"safeHTML", func(s string) template.HTML {
			return template.HTML(s)
		},
	)

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
		stic = static.New("./static", static.Config{
			FS:            srv.rt.Embeds["static"],
			Browse:        false,
			CacheDuration: 10 * time.Second, // TODO: Make configurable
		})
	}
	srv.app.Get("/static*", stic)

	if srv.r, err = root.New(srv.rt, srv.app); err != nil {
		return err
	}

	return nil
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
