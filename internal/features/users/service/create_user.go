package service

import (
	"context"
	"fmt"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

type CreateUserInput struct {
	FullName    string
	PhoneNumber *string
}

func (s *Service) CreateUser(ctx context.Context, in CreateUserInput) (domain.User, error) {
	user, err := domain.NewUser(in.FullName, in.PhoneNumber)
	if err != nil {
		return domain.User{}, err
	}

	created, err := s.repo.Create(ctx, user)
	if err != nil {
		return domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return created, nil
}
