package runtime

import (
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"fmt"
	"log/slog"
	"os"
	runt "runtime"
	"strings"

	"xn--gckvb8fzb.com/hyperuplink/services/config"
	"xn--gckvb8fzb.com/hyperuplink/services/database"
	"xn--gckvb8fzb.com/hyperuplink/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/services/intnat"
	"xn--gckvb8fzb.com/hyperuplink/services/magick"
	"xn--gckvb8fzb.com/hyperuplink/services/markdown"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories"
	"xn--gckvb8fzb.com/hyperuplink/services/storage"
)

var (
	Version string
	Commit  string
	Date    string
)

type Build struct {
	Version string
	Commit  string
	Date    string
	Hash    string
}

type Runtime struct {
	Build        Build
	Embeds       map[string]*embed.FS
	Config       *config.Config
	Logger       *slog.Logger
	LoggerLevel  slog.Level
	ALogger      AsyncLogger
	Database     *database.Database
	Repositories *repositories.Repositories
	Storage      *storage.Storage
	Magick       *magick.Magick
	Intnat       *intnat.Intnat
	Markdown     *markdown.Markdown
	Dispatch     *dispatch.Dispatch
}

const (
	ModeDevelopment string = "development"
	ModeProduction  string = "production"
)

func New(cfgstr string) (rt *Runtime, err error) {
	rt = new(Runtime)

	rt.Build.Version = Version
	rt.Build.Commit = Commit
	rt.Build.Date = Date
	rt.Build.Hash = rt.computeBuildHash(rt.Build.Commit, rt.Build.Date)

	rt.Embeds = make(map[string]*embed.FS)

	if rt.Config, err = config.New(cfgstr); err != nil {
		return nil, err
	}

	rt.LoggerLevel = slog.Level(0)
	if err = rt.LoggerLevel.UnmarshalText(rt.Config.LoggingLevel()); err != nil {
		return nil, err
	}

	rt.Logger = slog.New(
		slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: rt.LoggerLevel,
		}),
	)

	rt.ALogger = NewAsyncLogger(rt.Logger)

	rt.Debug("status", "exec")

	rt.Debug("new", "database")
	if rt.Database, err = database.New(
		rt.Logger,
		rt.Config.DatabaseConnection(),
	); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "repositories")
	if rt.Repositories, err = repositories.New(rt.Database, rt.Config); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "storagescfg")
	storagesCfg, err := rt.Config.Storages()
	if err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}
	rt.Debug("new", "storage")
	if rt.Storage, err = storage.New(storagesCfg); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "magick")
	if rt.Magick, err = magick.New(); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "intnat")
	if rt.Intnat, err = intnat.New(); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "markdown")
	if rt.Markdown, err = markdown.New(); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Debug("new", "rediscfg")
	redisCfg, err := rt.Config.Redis()
	if err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}
	rt.Debug("new", "dispatch")
	if rt.Dispatch, err = dispatch.New(redisCfg); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Info("status", "ok")
	return rt, err
}

func (rt *Runtime) Startup() (err error) {
	rt.Debug("status", "exec")

	rt.Debug("startup", "config")
	if err = rt.Config.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "database")
	if err = rt.Database.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "repositories")
	if err = rt.Repositories.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "storage")
	if err = rt.Storage.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "magick")
	if err = rt.Magick.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "intnat")
	if err = rt.Intnat.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "markdown")
	if err = rt.Markdown.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("startup", "dispatch")
	if err = rt.Dispatch.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Info("status", "ok")

	return nil
}

func (rt *Runtime) Shutdown() (err error) {
	rt.Debug("status", "exec")

	rt.Debug("shutdown", "dispatch")
	if err = rt.Dispatch.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "markdown")
	if err = rt.Markdown.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "intnat")
	if err = rt.Intnat.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "magick")
	if err = rt.Magick.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "storage")
	if err = rt.Storage.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "repositories")
	if err = rt.Repositories.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "database")
	if err = rt.Database.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Debug("shutdown", "config")
	if err = rt.Config.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Info("status", "ok")

	return nil
}

func (rt *Runtime) NilOrDie(err error) {
	if err != nil {
		fn := rt.getLogFnName()
		rt.Logger.Error(fn, "error", err)
		rt.Exit(1)
	}
}

func (rt *Runtime) Exit(code int) {
	rt.Shutdown()
	os.Exit(code)
}

func (rt *Runtime) getLogFnName() string {
	pc, _, _, ok := runt.Caller(2)
	if !ok {
		return "Unknown"
	}
	fn := runt.FuncForPC(pc)
	if fn == nil {
		return "Unknown"
	}
	fullName := fn.Name()
	fullSplit := strings.Split(fullName, ".")
	fSL := len(fullSplit)

	if fSL == 1 {
		return fullSplit[0]
	} else if fSL > 1 {
		pkg := fullSplit[fSL-2]
		if strings.Index(pkg, "/") > -1 {
			pkgs := strings.Split(pkg, "/")
			pkg = pkgs[len(pkgs)-1]
		}
		if strings.Index(pkg, "(") > -1 {
			pkg = strings.ReplaceAll(pkg, "(", "")
			pkg = strings.ReplaceAll(pkg, ")", "")
			pkg = strings.ReplaceAll(pkg, "*", "")
		}
		mtd := fullSplit[fSL-1]
		return fmt.Sprintf("%s.%s", pkg, mtd)
	}

	return "Unknown"
}

func (rt *Runtime) Debug(args ...any) {
	fn := rt.getLogFnName()
	rt.Logger.Debug(fn, args...)
}

func (rt *Runtime) Info(args ...any) {
	fn := rt.getLogFnName()
	rt.Logger.Info(fn, args...)
}

func (rt *Runtime) Warn(args ...any) {
	fn := rt.getLogFnName()
	rt.Logger.Warn(fn, args...)
}

func (rt *Runtime) Error(args ...any) {
	fn := rt.getLogFnName()
	rt.Logger.Error(fn, args...)
}

func (rt *Runtime) computeBuildHash(args ...string) string {
	h := sha256.New()
	h.Write([]byte(strings.Join(args, "")))
	hashBytes := h.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func (rt *Runtime) IsDevelopmentMode() bool {
	if rt.Config.GeneralMode() == ModeDevelopment {
		return true
	}

	return false
}
