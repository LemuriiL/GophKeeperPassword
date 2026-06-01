package main

import (
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

	app, err := client.NewApp()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
