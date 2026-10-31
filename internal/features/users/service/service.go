package service

import (
	"context"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user domain.User) (domain.User, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

type Service struct {
	repo UserRepository
}

func New(repo UserRepository) *Service {
	return &Service{repo: repo}
}
