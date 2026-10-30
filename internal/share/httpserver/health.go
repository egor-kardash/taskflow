package httpserver

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type statusResponse struct {
	Status string `json:"status"`
}

type Pinger interface {
	Ping(ctx context.Context) error
}

const readinessTimeout = 2 * time.Second

func healthz(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, statusResponse{Status: "ok"})
}

func readyz(log *slog.Logger, db Pinger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), readinessTimeout)
		defer cancel()

		if err := db.Ping(ctx); err != nil {
			log.Debug("readiness check failed", "error", err)
			WriteError(w, http.StatusServiceUnavailable, "database_unavailable", "database is not reachable")
			return
		}

		WriteJSON(w, http.StatusOK, statusResponse{Status: "ready"})
	}
}
