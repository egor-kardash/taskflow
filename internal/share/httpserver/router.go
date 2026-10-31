package httpserver

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(log *slog.Logger, db Pinger) chi.Router {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(RequestLogger(log))
	r.Use(Recoverer(log))

	r.NotFound(func(w http.ResponseWriter, _ *http.Request) {
		WriteError(w, http.StatusNotFound, "not_found", "resource not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, _ *http.Request) {
		WriteError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
	})

	r.Get("/healthz", healthz)
	r.Get("/readyz", readyz(log, db))

	return r
}
