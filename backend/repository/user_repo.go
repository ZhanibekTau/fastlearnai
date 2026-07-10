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
