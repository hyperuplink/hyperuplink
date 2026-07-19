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

	"xn--gckvb8fzb.com/hyperuplink/cron"
	"xn--gckvb8fzb.com/hyperuplink/http"
	logicactivity "xn--gckvb8fzb.com/hyperuplink/logic/helpers/activity"
	logicsession "xn--gckvb8fzb.com/hyperuplink/logic/root/session"
	"xn--gckvb8fzb.com/hyperuplink/runtime"
	"xn--gckvb8fzb.com/hyperuplink/worker"
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
	var web *http.HTTP
	var api *http.HTTP
	var wrk *worker.Worker
	var crn *cron.Cron
	var err error

	flag.Parse()

	if flagVersion {
		fmt.Printf("Hyper Uplink %s\nCommit: %s\nBuild date: %s\n",
			runtime.Version,
			runtime.Commit,
			runtime.Date,
		)
		os.Exit(0)
	}

	rt, err = runtime.New(flagCfgstr)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

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

	if flagCreateUser != "" {
		rt.Embeds["migrations"] = &embedMigrations
		rt.Database.SetMigrations(rt.Embeds["migrations"])
		if err = createUser(rt, flagCreateUser); err != nil {
			fmt.Printf("%s\n", err)
			os.Exit(1)
		}
		os.Exit(0)
	}

	rt.Embeds["migrations"] = &embedMigrations
	rt.Embeds["static"] = &embedStaticFiles
	rt.Embeds["views"] = &embedViews
	rt.Embeds["locales"] = &embedLocales
	rt.Embeds["templates"] = &embedTemplates
	rt.Embeds["docs"] = &embedDocs

	rt.Database.SetMigrations(rt.Embeds["migrations"])
	rt.Intnat.SetLocales(rt.Embeds["locales"])

	err = rt.Startup()
	rt.NilOrDie(err)

	// ---[ WEB ]-------------------------------------------------------------- //
	web, err = http.New(rt, http.IfaceWeb)
	rt.NilOrDie(err)

	err = web.Startup()
	rt.NilOrDie(err)

	go web.Run()

	// ---[ API ]-------------------------------------------------------------- //
	api, err = http.New(rt, http.IfaceAPI)
	rt.NilOrDie(err)

	err = api.Startup()
	rt.NilOrDie(err)

	go api.Run()

	// ---[ WORKER ]----------------------------------------------------------- //
	wrk, err = worker.New(rt)
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
	api.Shutdown()
	web.Shutdown()
	wrk.Shutdown()

	rt.Exit(0)
}

func registerCronFunctions(rt *runtime.Runtime) (err error) {
	return rt.Cron.Register(
		logicactivity.CleanupAdminLogID,
		logicactivity.CleanupAdminLogSpec,
		func() error { return logicactivity.CleanupAdminLog(rt) },
	)
}

func createUser(rt *runtime.Runtime, jsonStr string) error {
	if err := rt.Database.Startup(); err != nil {
		return err
	}
	defer rt.Database.Shutdown()

	if err := rt.Repositories.Startup(); err != nil {
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
