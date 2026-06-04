package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log/slog"
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

var loadServerConfig = server.LoadConfig

var newServerApp = func(cfg server.Config) (serverApp, error) {
	return server.NewApp(cfg)
}

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout io.Writer, stderr io.Writer) int {
	buildinfo.Print(buildVersion, buildDate, buildCommit)

	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	fs.SetOutput(stderr)

	cfgPathShort := fs.String("c", "", "path to config")
	cfgPathLong := fs.String("config", "", "path to config")

	if err := fs.Parse(args); err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	cfg, err := loadServerConfig(*cfgPathShort, *cfgPathLong)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}

	app, err := newServerApp(cfg)
	if err != nil {
		fmt.Fprintln(stderr, err)
		return 1
	}
	defer func() {
		if err := app.Close(); err != nil {
			slog.Error("close app", "error", err)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	select {
	case err := <-errCh:
		if err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintln(stderr, err)
			return 1
		}
	}

	return 0
}
