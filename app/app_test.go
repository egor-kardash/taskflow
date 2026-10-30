package app

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/egor-kardash/taskflow/internal/config"
	"github.com/egor-kardash/taskflow/internal/share/httpserver"
)

func testApp(t *testing.T, addr string) *App {
	t.Helper()

	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	srv := httpserver.New(config.HTTPServer{
		Address:     addr,
		Timeout:     time.Second,
		IdleTimeout: time.Second,
	}, http.NotFoundHandler(), log)

	return &App{
		server:          srv,
		log:             log,
		shutdownTimeout: time.Second,
	}
}

func TestRunStopsOnContextCancel(t *testing.T) {
	a := testApp(t, "127.0.0.1:0")

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- a.Run(ctx) }()

	cancel()

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run did not return after cancel")
	}
}

func TestRunReturnsListenError(t *testing.T) {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	defer func() { _ = ln.Close() }()

	a := testApp(t, ln.Addr().String())

	if err := a.Run(context.Background()); err == nil {
		t.Fatal("expected a listen error for an occupied address")
	}
}
