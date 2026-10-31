package service

import (
	"context"
	"fmt"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

func (s *Service) GetUser(ctx context.Context, id int64) (domain.User, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return domain.User{}, fmt.Errorf("get user %d: %w", id, err)
	}

	return user, nil
}
