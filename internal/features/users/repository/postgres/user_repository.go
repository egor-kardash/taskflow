package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/egor-kardash/taskflow/internal/features/users/domain"
)

const (
	uniqueViolation = "23505"
	phoneNumberKey  = "users_phone_number_key"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{pool: pool}
}

func (r *UserRepository) Create(ctx context.Context, user domain.User) (domain.User, error) {
	const query = `
		INSERT INTO users (full_name, phone_number)
		VALUES ($1, $2)
		RETURNING id, version, full_name, phone_number`

	var created domain.User
	err := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber).Scan(
		&created.ID,
		&created.Version,
		&created.FullName,
		&created.PhoneNumber,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == uniqueViolation && pgErr.ConstraintName == phoneNumberKey {
			return domain.User{}, fmt.Errorf("insert user: %w", domain.ErrPhoneTaken)
		}
		return domain.User{}, fmt.Errorf("insert user: %w", err)
	}

	return created, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int64) (domain.User, error) {
	const query = `
		SELECT id, version, full_name, phone_number
		FROM users
		WHERE id = $1`

	var user domain.User
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Version,
		&user.FullName,
		&user.PhoneNumber,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.User{}, fmt.Errorf("select user %d: %w", id, domain.ErrNotFound)
		}
		return domain.User{}, fmt.Errorf("select user %d: %w", id, err)
	}

	return user, nil
}
