package service

import (
	"context"
	"errors"
	"testing"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

type fakeRepo struct {
	created   domain.User
	createErr error
	user      domain.User
	getErr    error

	lastUser domain.User
	lastID   int64
}

func (f *fakeRepo) Create(_ context.Context, user domain.User) (domain.User, error) {
	f.lastUser = user
	return f.created, f.createErr
}

func (f *fakeRepo) GetByID(_ context.Context, id int64) (domain.User, error) {
	f.lastID = id
	return f.user, f.getErr
}

func TestCreateUserPassesDomainUser(t *testing.T) {
	repo := &fakeRepo{created: domain.User{ID: 1, Version: 1, FullName: "Egor"}}
	phone := "  +79991234567  "

	user, err := New(repo).CreateUser(context.Background(), CreateUserInput{
		FullName:    "  Egor  ",
		PhoneNumber: &phone,
	})
	if err != nil {
		t.Fatalf("CreateUser: %v", err)
	}
	if user.ID != 1 {
		t.Fatalf("id = %d, want 1", user.ID)
	}
	if repo.lastUser.FullName != "Egor" {
		t.Fatalf("name = %q", repo.lastUser.FullName)
	}
	if repo.lastUser.PhoneNumber == nil || *repo.lastUser.PhoneNumber != "+79991234567" {
		t.Fatalf("phone = %v", repo.lastUser.PhoneNumber)
	}
}

func TestCreateUserDoesNotCallRepoOnValidationError(t *testing.T) {
	repo := &fakeRepo{}

	_, err := New(repo).CreateUser(context.Background(), CreateUserInput{FullName: "  "})

	var validationErr *domain.ValidationError
	if !errors.As(err, &validationErr) {
		t.Fatalf("err = %v, want ValidationError", err)
	}
	if repo.lastUser.FullName != "" || repo.lastUser.PhoneNumber != nil {
		t.Fatalf("repo was called: %+v", repo.lastUser)
	}
}

func TestCreateUserWrapsRepoError(t *testing.T) {
	repo := &fakeRepo{createErr: domain.ErrPhoneTaken}

	_, err := New(repo).CreateUser(context.Background(), CreateUserInput{FullName: "Egor"})
	if !errors.Is(err, domain.ErrPhoneTaken) {
		t.Fatalf("err = %v", err)
	}
}
