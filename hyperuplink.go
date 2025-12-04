package main

import (
	"embed"
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mrusme/hyperuplink/http"
	"github.com/mrusme/hyperuplink/runtime"
	"github.com/mrusme/hyperuplink/worker"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

//go:embed static/*
var embedStaticFiles embed.FS

//go:embed views/*
var embedViews embed.FS

//go:embed locales/*
var embedLocales embed.FS

var (
	flagCfgstr  string
	flagVersion bool
)

func init() {
	flag.StringVar(&flagCfgstr, "c", "file:///etc/hyperuplink.toml", "configuration string")
	flag.BoolVar(&flagVersion, "v", false, "Print version information and exit")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Use: %s [-opts]\n\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var err error

	flag.Parse()

	if flagVersion {
		fmt.Printf("Hyper Uplink %s\nCommit: %s\nBuild date: %s\n",
			runtime.Version,
			runtime.Commit,
			runtime.Version,
		)
	}

	rt, err := runtime.New(flagCfgstr)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	rt.Embeds["migrations"] = &embedMigrations
	rt.Embeds["static"] = &embedStaticFiles
	rt.Embeds["views"] = &embedViews
	rt.Embeds["locales"] = &embedLocales

	rt.Database.SetMigrations(rt.Embeds["migrations"])
	rt.Intnat.SetLocales(rt.Embeds["locales"])

	err = rt.Startup()
	rt.NilOrDie(err)

	// ---[ WEB ]-------------------------------------------------------------- //
	web, err := http.New(rt, http.IfaceWeb)
	rt.NilOrDie(err)

	err = web.Startup()
	rt.NilOrDie(err)

	go web.Run()

	// ---[ WORKER ]----------------------------------------------------------- //
	wrk, err := worker.New(rt)
	rt.NilOrDie(err)

	err = wrk.Startup()
	rt.NilOrDie(err)

	go wrk.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	web.Shutdown()

	rt.Exit(0)
}
