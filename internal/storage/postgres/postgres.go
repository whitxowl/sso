package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"sso/internal/domain/models"
	"sso/internal/storage"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Storage struct {
	db *sql.DB
}

func New(dsn string) (*Storage, error) {
	const op = "storage.New"

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

// SaveUser saves user to db
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.SaveUser"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "INSERT INTO users (email, pass_hash) VALUES ($1, $2)"

	res, err := s.db.ExecContext(ctx, query, email, passHash)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.User"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT * FROM users WHERE email = $1"

	user := models.User{}

	err := s.db.QueryRowContext(ctx, query, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return user, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return user, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "storage.IsAdmin"

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	query := "SELECT is_admin FROM users WHERE id = $1"

	var isAdmin bool

	err := s.db.QueryRowContext(ctx, query, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
