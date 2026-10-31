package http

import (
	"context"
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
	"github.com/egor-kardash/taskflow/internal/features/users/service"
)

const maxBodyBytes = 1 << 20 // 1 MiB

type UsersService interface {
	CreateUser(ctx context.Context, in service.CreateUserInput) (domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
}

type Handler struct {
	users UsersService
	log   *slog.Logger
}

func New(users UsersService, log *slog.Logger) *Handler {
	return &Handler{users: users, log: log}
}

func (h *Handler) Register(r chi.Router) {
	r.Post("/users", h.createUser)
	r.Get("/users/{id}", h.getUser)
}
