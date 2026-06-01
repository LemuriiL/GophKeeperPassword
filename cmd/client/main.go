package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/LemuriiL/GophKeeperPassword/internal/buildinfo"
	"github.com/LemuriiL/GophKeeperPassword/internal/client"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	buildinfo.Print(buildVersion, buildDate, buildCommit)

	fs := flag.NewFlagSet("client", flag.ContinueOnError)
	configShort := fs.String("c", "", "path to config")
	configLong := fs.String("config", "", "path to config")

	if err := fs.Parse(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	app, err := client.NewApp(*configShort, *configLong)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := app.Run(fs.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
