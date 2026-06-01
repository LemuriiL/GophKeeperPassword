package main

import (
	"context"
	"flag"
	"log"
	"log/slog"
	"os/signal"
	"syscall"
	"time"

	"github.com/LemuriiL/GophKeeperPassword/internal/buildinfo"
	"github.com/LemuriiL/GophKeeperPassword/internal/server"
)

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

func main() {
	buildinfo.Print(buildVersion, buildDate, buildCommit)

	cfgPathShort := flag.String("c", "", "path to config")
	cfgPathLong := flag.String("config", "", "path to config")
	flag.Parse()

	cfg, err := server.LoadConfig(*cfgPathShort, *cfgPathLong)
	if err != nil {
		log.Fatal(err)
	}

	app, err := server.NewApp(cfg)
	if err != nil {
		log.Fatal(err)
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
			log.Fatal(err)
		}
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		if err := app.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
	}
}
