package repository

import (
	"context"
	"errors"
	"time"

	"github.com/fastlearn/backend/models"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type CreateUserInput struct {
	Email        string
	PasswordHash string
	Name         string
	OAuthProvider *string
	OAuthID      *string
}

type UserRepository interface {
	ExistsByEmail(ctx context.Context, email string) (bool, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByID(ctx context.Context, id string) (*models.User, error)
	Create(ctx context.Context, input CreateUserInput) (*models.User, error)
	SearchByEmail(ctx context.Context, query string) ([]models.User, error)
	UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error
	Count(ctx context.Context) (int, error)
}

type pgUserRepo struct{ db *pgxpool.Pool }

func NewUserRepo(db *pgxpool.Pool) UserRepository {
	return &pgUserRepo{db: db}
}

func (r *pgUserRepo) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var exists bool
	err := r.db.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM users WHERE email=$1)`, email,
	).Scan(&exists)
	return exists, err
}

func (r *pgUserRepo) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, name, password_hash, oauth_provider, created_at, updated_at
		 FROM users WHERE email=$1`, email,
	).Scan(&u.ID, &u.Email, &u.Name, &u.PasswordHash, &u.OAuthProvider, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

func (r *pgUserRepo) FindByID(ctx context.Context, id string) (*models.User, error) {
	var u models.User
	err := r.db.QueryRow(ctx,
		`SELECT id, email, name, oauth_provider, created_at, updated_at
		 FROM users WHERE id=$1`, id,
	).Scan(&u.ID, &u.Email, &u.Name, &u.OAuthProvider, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return &u, err
}

// SearchByEmail does a case-insensitive partial match, for admin lookup.
func (r *pgUserRepo) SearchByEmail(ctx context.Context, query string) ([]models.User, error) {
	rows, err := r.db.Query(ctx,
		`SELECT id, email, name, oauth_provider, created_at, updated_at
		 FROM users WHERE email ILIKE $1 ORDER BY created_at DESC LIMIT 20`,
		"%"+query+"%",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.ID, &u.Email, &u.Name, &u.OAuthProvider, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

func (r *pgUserRepo) Count(ctx context.Context) (int, error) {
	var count int
	err := r.db.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&count)
	return count, err
}

func (r *pgUserRepo) UpdatePasswordHash(ctx context.Context, userID, passwordHash string) error {
	_, err := r.db.Exec(ctx,
		`UPDATE users SET password_hash=$1, updated_at=now() WHERE id=$2`,
		passwordHash, userID,
	)
	return err
}

func (r *pgUserRepo) Create(ctx context.Context, input CreateUserInput) (*models.User, error) {
	id := uuid.New().String()
	now := time.Now()

	_, err := r.db.Exec(ctx,
		`INSERT INTO users (id, email, password_hash, name, oauth_provider, oauth_id, created_at, updated_at)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$7)`,
		id, input.Email, input.PasswordHash, input.Name, input.OAuthProvider, input.OAuthID, now,
	)
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           id,
		Email:        input.Email,
		Name:         &input.Name,
		OAuthProvider: input.OAuthProvider,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}
