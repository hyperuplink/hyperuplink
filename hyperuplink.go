package main

import (
	"context"
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
	servicescron "xn--gckvb8fzb.com/glides/services/cron"
	"xn--gckvb8fzb.com/glides/services/database"
	glidesdispatch "xn--gckvb8fzb.com/glides/services/dispatch"
	"xn--gckvb8fzb.com/glides/services/intnat"
	"xn--gckvb8fzb.com/glides/services/markdown"
	"xn--gckvb8fzb.com/glides/services/storage"

	"xn--gckvb8fzb.com/glides/cron"
	glideshttp "xn--gckvb8fzb.com/glides/http"
	"xn--gckvb8fzb.com/glides/runtime"
	"xn--gckvb8fzb.com/glides/worker"
	gh "xn--gckvb8fzb.com/hyperuplink/helpers"
	"xn--gckvb8fzb.com/hyperuplink/http/api"
	"xn--gckvb8fzb.com/hyperuplink/http/routes"
	"xn--gckvb8fzb.com/hyperuplink/http/web"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/services/activity"
	"xn--gckvb8fzb.com/hyperuplink/services/dispatch"
	"xn--gckvb8fzb.com/hyperuplink/services/magick"
	"xn--gckvb8fzb.com/hyperuplink/services/repositories"
	"xn--gckvb8fzb.com/hyperuplink/worker/targets"
)

var (
	Version string
	Commit  string
	Date    string
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

//go:embed static/*
var embedStaticFiles embed.FS

//go:embed views/*
var embedViews embed.FS

//go:embed locales/*
var embedLocales embed.FS

//go:embed templates/*
var embedTemplates embed.FS

//go:embed docs/*
var embedDocs embed.FS

var (
	flagCfgstr     string
	flagVersion    bool
	flagReset      string
	flagCreateUser string
)

func init() {
	flag.StringVar(&flagCfgstr, "c", "file:///etc/hyperuplink.toml", "configuration string")
	flag.BoolVar(&flagVersion, "v", false, "Print version information and exit")
	flag.StringVar(&flagReset, "reset", "", "Clear the whole database and exit (requires the current time as HH:MM (24h) confirmation, e.g. --reset 10:42)")
	flag.StringVar(&flagCreateUser, "create-user", "", `Create an activated user from a JSON object of signup fields and exit, e.g. --create-user '{"username":"dummy1","email":"dummy1@example.com","password":"mypassword"}'`)
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use: %s [-opts]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var rt *runtime.Runtime
	var websrv *glideshttp.HTTP
	var apisrv *glideshttp.HTTP
	var wrk *worker.Worker
	var crn *cron.Cron
	var err error

	flag.Parse()

	if flagVersion {
		fmt.Printf("Hyper Uplink %s\nCommit: %s\nBuild date: %s\n",
			Version,
			Commit,
			Date,
		)
		os.Exit(0)
	}

	if rt, err = newRuntime(flagCfgstr); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if flagReset != "" {
		if err = runtime.ValidateResetTime(flagReset, time.Now()); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		if err = database.From(rt).Reset(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	rt.AddEmbed("migrations", &embedMigrations)
	database.From(rt).SetMigrations(rt.GetEmbed("migrations"))

	if flagCreateUser != "" {
		if err = createUser(rt, flagCreateUser); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	err = prepare(rt)
	rt.NilOrDie(err)

	err = rt.Startup()
	rt.NilOrDie(err)

	routes.Use()

	webcfg, err := gh.WebConfigFrom(rt.Config())
	rt.NilOrDie(err)
	apicfg, err := gh.APIConfigFrom(rt.Config())
	rt.NilOrDie(err)

	wrk, err = worker.New(rt, targets.Register)
	rt.NilOrDie(err)
	processes := []runtime.IProcess{wrk}

	if webcfg.Enable {
		webiface, err := web.New(rt, webcfg)
		rt.NilOrDie(err)
		websrv, err = glideshttp.New(rt, webiface)
		rt.NilOrDie(err)
		processes = append(processes, websrv)
	}

	if apicfg.Enable {
		apiiface, err := api.New(rt, apicfg)
		rt.NilOrDie(err)
		apisrv, err = glideshttp.New(rt, apiiface)
		rt.NilOrDie(err)
		processes = append(processes, apisrv)
	}

	crn, err = cron.New(rt)
	rt.NilOrDie(err)
	err = registerCronFunctions(rt)
	rt.NilOrDie(err)
	processes = append(processes, crn)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if err = rt.Serve(ctx, processes...); err != nil {
		rt.Exit(1)
	}

	rt.Exit(0)
}

func newRuntime(cfgstr string) (rt *runtime.Runtime, err error) {
	if rt, err = runtime.New(runtime.Opts{
		Cfgstr:  cfgstr,
		Version: Version,
		Commit:  Commit,
		Date:    Date,
	}); err != nil {
		return nil, err
	}

	if _, err = database.Register(rt); err != nil {
		return nil, err
	}
	if _, err = storage.Register(rt); err != nil {
		return nil, err
	}
	if _, err = intnat.Register(rt); err != nil {
		return nil, err
	}
	if _, err = markdown.Register(rt); err != nil {
		return nil, err
	}
	if _, err = glidesdispatch.Register(rt); err != nil {
		return nil, err
	}
	if _, err = servicescron.Register(rt); err != nil {
		return nil, err
	}

	return rt, nil
}

func prepare(rt *runtime.Runtime) (err error) {
	rt.AddEmbed("static", &embedStaticFiles)
	rt.AddEmbed("views", &embedViews)
	rt.AddEmbed("locales", &embedLocales)
	rt.AddEmbed("templates", &embedTemplates)
	rt.AddEmbed("docs", &embedDocs)

	intnat.From(rt).SetLocales(rt.GetEmbed("locales"), gh.DefaultLocale, gh.Locales...)

	var srv any
	rt.Debug("new", "repositories")
	if srv, err = repositories.New(database.From(rt), rt.Config()); err != nil {
		return err
	}
	rt.AddService("repositories", srv)

	rt.Debug("new", "activity")
	if srv, err = activity.New(rt.Logger(), gh.Repositories(rt).Activity); err != nil {
		return err
	}
	rt.AddService("activity", srv)

	rt.Debug("new", "magick")
	if srv, err = magick.New(); err != nil {
		return err
	}
	rt.AddService("magick", srv)

	rt.Debug("new", "dispatch")
	if srv, err = dispatch.New(
		glidesdispatch.From(rt), gh.Repositories(rt).Setting,
	); err != nil {
		return err
	}
	rt.AddService("dispatch", srv)

	rt.OnStartup(
		gh.Repositories(rt).Startup,
		gh.Activity(rt).Startup,
		gh.Magick(rt).Startup,
	)
	rt.OnShutdown(
		gh.Magick(rt).Shutdown,
		gh.Activity(rt).Shutdown,
		gh.Repositories(rt).Shutdown,
	)

	return nil
}

func registerCronFunctions(rt *runtime.Runtime) (err error) {
	return servicescron.From(rt).Register(
		logicactivity.CleanupAdminLogID,
		logicactivity.CleanupAdminLogSpec,
		func(ctx context.Context) error { return logicactivity.CleanupAdminLog(ctx, rt) },
	)
}

func createUser(rt *runtime.Runtime, jsonStr string) error {
	if err := prepare(rt); err != nil {
		return err
	}

	if err := intnat.From(rt).Startup(); err != nil {
		return err
	}

	if err := database.From(rt).Startup(); err != nil {
		return err
	}
	defer database.From(rt).Shutdown()

	if err := gh.Repositories(rt).Startup(); err != nil {
		return err
	}

	in := new(logicsession.SignUpInput)
	if err := json.Unmarshal([]byte(jsonStr), in); err != nil {
		return err
	}

	usr, err := logicsession.SignUp(context.Background(), rt, in, logicsession.SignUpOptions{Activate: true})
	if err != nil {
		return err
	}

	rt.Info("create-user", "ok",
		"id", usr.ID.String(),
		"username", usr.Username,
		"role", string(usr.Role),
	)
	return nil
}
