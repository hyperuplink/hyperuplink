package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/mrusme/hyperuplink/http"
	"github.com/mrusme/hyperuplink/runtime"
)

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

	err = rt.Startup()
	rt.NilOrDie(err)

	web, err := http.New(rt, http.IfaceWeb)
	rt.NilOrDie(err)

	err = web.Startup()
	rt.NilOrDie(err)

	go web.Run()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	web.Shutdown()

	rt.Exit(0)
}
