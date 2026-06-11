package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/LemuriiL/GophKeeperPassword/internal/buildinfo"
	"github.com/LemuriiL/GophKeeperPassword/internal/client"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

type clientApp interface {
	Run(args []string) error
}

var newClientApp = func(configShort string, configLong string) (clientApp, error) {
	return client.NewApp(configShort, configLong)
}

// main запускает клиент
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run выполняет клиентскую команду
func run(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("client", flag.ContinueOnError)
	fs.SetOutput(stderr)

	configShort := fs.String("c", "", "path to config")
	configLong := fs.String("config", "", "path to config")
	version := fs.Bool("version", false, "print build info")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if *version || len(fs.Args()) > 0 && fs.Args()[0] == "version" {
		buildinfo.Print(stdout, buildVersion, buildDate, buildCommit)
		return 0
	}

	app, err := newClientApp(*configShort, *configLong)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	if err := app.Run(fs.Args()); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	return 0
}
