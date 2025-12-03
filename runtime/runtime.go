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

	"github.com/mrusme/hyperuplink/services/config"
	"github.com/mrusme/hyperuplink/services/database"
	"github.com/mrusme/hyperuplink/services/intnat"
	"github.com/mrusme/hyperuplink/services/repositories"
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
	Database     *database.Database
	Repositories *repositories.Repositories
	Intnat       *intnat.Intnat
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

	rt.Debug("status", "exec")

	if rt.Database, err = database.New(
		rt.Logger,
		rt.Config.DatabaseConnection(),
	); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	if rt.Repositories, err = repositories.New(rt.Database); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	if rt.Intnat, err = intnat.New(); err != nil {
		rt.Error("status", "error", "error", err)
		return nil, err
	}

	rt.Info("status", "ok")
	return rt, err
}

func (rt *Runtime) Startup() (err error) {
	rt.Debug("status", "exec")

	if err = rt.Config.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	if err = rt.Database.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	if err = rt.Intnat.Startup(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	rt.Info("status", "ok")

	return nil
}

func (rt *Runtime) Shutdown() (err error) {
	rt.Debug("status", "exec")

	if err = rt.Intnat.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

	if err = rt.Database.Shutdown(); err != nil {
		rt.Error("status", "error", "error", err)
		return err
	}

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
