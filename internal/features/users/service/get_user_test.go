package service

import (
	"context"
	"errors"
	"testing"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

func TestGetUserPassesID(t *testing.T) {
	repo := &fakeRepo{user: domain.User{ID: 42, Version: 1, FullName: "Egor"}}

	user, err := New(repo).GetUser(context.Background(), 42)
	if err != nil {
		t.Fatalf("GetUser: %v", err)
	}
	if repo.lastID != 42 || user.ID != 42 {
		t.Fatalf("lastID = %d, id = %d", repo.lastID, user.ID)
	}
}

func TestGetUserWrapsNotFound(t *testing.T) {
	repo := &fakeRepo{getErr: domain.ErrNotFound}

	_, err := New(repo).GetUser(context.Background(), 999)
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("err = %v", err)
	}
}
