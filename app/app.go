package app

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/egor-kardash/taskflow/internal/config"
	"github.com/egor-kardash/taskflow/internal/share/postgres"
	"github.com/jackc/pgx/v5/pgxpool"
)

type App struct {
	pool *pgxpool.Pool
}

func New(ctx context.Context, cfg *config.Config) (*App, error) {
	pool, err := postgres.NewPool(ctx, cfg.Database.DSN())
	if err != nil {
		return nil, fmt.Errorf("init postgres: %w", err)
	}

	return &App{pool: pool}, nil
}

func (a *App) Run(ctx context.Context) error {
	<-ctx.Done()
	slog.Info("shutting down")
	return nil
}

func (a *App) Close() error {
	a.pool.Close()
	return nil
}
