package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LemuriiL/GophKeeperPassword/internal/buildinfo"
	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

type serverApp interface {
	Run() error
	Shutdown(ctx context.Context) error
	Close() error
}

var newServerApp = func(cfg server.Config) (serverApp, error) {
	return server.NewApp(cfg)
}

// main запускает сервер
func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

// run выполняет запуск сервера
func run(args []string, stdout io.Writer, stderr io.Writer) int {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
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

	cfg, err := server.LoadConfig(*configShort, *configLong)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	app, err := newServerApp(cfg)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer app.Close()

	errCh := make(chan error, 1)

	go func() {
		errCh <- app.Run()
	}()

	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer signal.Stop(stopCh)

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	case <-stopCh:
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err = app.Shutdown(ctx); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	return 0
}
