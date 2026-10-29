package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/egor-kardash/taskflow/app"
	"github.com/egor-kardash/taskflow/internal/config"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	var configPath string
	flag.StringVar(&configPath, "config", os.Getenv("CONFIG_PATH"), "path to config file")
	flag.Parse()

	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	startCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	a, err := app.New(startCtx, cfg)
	if err != nil {
		return err
	}
	defer func() {
		err = errors.Join(err, a.Close())
	}()

	slog.Info("connected to postgres",
		"host", cfg.Database.Host,
		"port", cfg.Database.Port,
		"db", cfg.Database.Name,
	)
	slog.Info("app started, press Ctrl+C to stop")

	return a.Run(ctx)
}
