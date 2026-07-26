package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	if rt, err = runtime.New(runtime.Opts{
		Cfgstr:  flagCfgstr,
		Version: Version,
		Commit:  Commit,
		Date:    Date,
		Services: runtime.Services{
			Database: true,
			Storage:  true,
			Intnat:   true,
			Markdown: true,
			Dispatch: true,
			Cron:     true,
		},
	}); err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	rt.SetBuild(Version, Commit, Date)

	if flagReset != "" {
		if err = runtime.ValidateResetTime(flagReset, time.Now()); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		if err = rt.Reset(); err != nil {
			os.Exit(1)
		}
		os.Exit(0)
	}

	rt.AddEmbed("migrations", &embedMigrations)
	rt.Database().SetMigrations(rt.GetEmbed("migrations"))

	if flagCreateUser != "" {
		if err = createUser(rt, flagCreateUser); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	rt.AddEmbed("static", &embedStaticFiles)
	rt.AddEmbed("views", &embedViews)
	rt.AddEmbed("locales", &embedLocales)
	rt.AddEmbed("templates", &embedTemplates)
	rt.AddEmbed("docs", &embedDocs)

	rt.Intnat().SetLocales(rt.GetEmbed("locales"), "en",
		"de", "en", "es", "fr", "it", "ro")

	var srv any
	rt.Debug("new", "repositories")
	if srv, err = repositories.New(rt.Database(), rt.Config()); err != nil {
		rt.Error("status", "error", "error", err)
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	rt.AddService("repositories", srv)

	rt.Debug("new", "activity")
	if srv, err = activity.New(rt.Logger(), gh.Repositories(rt).Activity); err != nil {
		rt.Error("status", "error", "error", err)
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	rt.AddService("activity", srv)

	rt.Debug("new", "magick")
	if srv, err = magick.New(); err != nil {
		rt.Error("status", "error", "error", err)
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
	rt.AddService("magick", srv)

	rt.Debug("new", "dispatch")
	if srv, err = dispatch.New(
		rt.Dispatch(), gh.Repositories(rt).Setting,
	); err != nil {
		rt.Error("status", "error", "error", err)
		fmt.Printf("%s\n", err)
		os.Exit(1)
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

	err = rt.Startup()
	rt.NilOrDie(err)

	routes.Use()

	// ---[ WEB ]-------------------------------------------------------------- //
	if rt.Config().WebEnable() {
		webiface, err := web.New(rt)
		rt.NilOrDie(err)

		websrv, err = glideshttp.New(rt, webiface)
		rt.NilOrDie(err)

		err = websrv.Startup()
		rt.NilOrDie(err)

		go websrv.Run()
	}

	// ---[ API ]-------------------------------------------------------------- //
	if rt.Config().APIEnable() {
		apiiface, err := api.New(rt)
		rt.NilOrDie(err)

		apisrv, err = glideshttp.New(rt, apiiface)
		rt.NilOrDie(err)

		err = apisrv.Startup()
		rt.NilOrDie(err)

		go apisrv.Run()
	}

	// ---[ WORKER ]----------------------------------------------------------- //
	wrk, err = worker.New(rt, targets.Register)
	rt.NilOrDie(err)

	err = wrk.Startup()
	rt.NilOrDie(err)

	go wrk.Run()

	// ---[ CRON ]------------------------------------------------------------- //
	crn, err = cron.New(rt)
	rt.NilOrDie(err)

	err = crn.Startup()
	rt.NilOrDie(err)

	err = registerCronFunctions(rt)
	rt.NilOrDie(err)

	go crn.Run()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	crn.Shutdown()
	if apisrv != nil {
		apisrv.Shutdown()
	}
	if websrv != nil {
		websrv.Shutdown()
	}
	wrk.Shutdown()

	rt.Exit(0)
}

func registerCronFunctions(rt *runtime.Runtime) (err error) {
	return rt.Cron().Register(
		logicactivity.CleanupAdminLogID,
		logicactivity.CleanupAdminLogSpec,
		func() error { return logicactivity.CleanupAdminLog(rt) },
	)
}

func createUser(rt *runtime.Runtime, jsonStr string) error {
	if err := rt.Database().Startup(); err != nil {
		return err
	}
	defer rt.Database().Shutdown()

	if err := gh.Repositories(rt).Startup(); err != nil {
		return err
	}

	in := new(logicsession.SignUpInput)
	if err := json.Unmarshal([]byte(jsonStr), in); err != nil {
		return err
	}

	usr, err := logicsession.SignUp(rt, in, logicsession.SignUpOptions{Activate: true})
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
