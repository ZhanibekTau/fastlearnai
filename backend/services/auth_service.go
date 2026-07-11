package services

import (
	"context"
	"errors"
	"log/slog"

	"github.com/fastlearn/backend/config"
	"github.com/fastlearn/backend/models"
	"github.com/fastlearn/backend/repository"
	"github.com/fastlearn/backend/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailExists        = errors.New("email already exists")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthService struct {
	users repository.UserRepository
	subs  repository.SubscriptionRepository
	email *EmailService
	cfg   *config.Config
}

func NewAuthService(
	users repository.UserRepository,
	subs repository.SubscriptionRepository,
	email *EmailService,
	cfg *config.Config,
) *AuthService {
	return &AuthService{users: users, subs: subs, email: email, cfg: cfg}
}

func (s *AuthService) Register(ctx context.Context, req models.RegisterRequest) (*models.AuthResponse, error) {
	exists, err := s.users.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrEmailExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user, err := s.users.Create(ctx, repository.CreateUserInput{
		Email:        req.Email,
		PasswordHash: string(hash),
		Name:         req.Name,
	})
	if err != nil {
		return nil, err
	}

	// Provision free subscription
	_ = s.subs.CreateDefault(ctx, user.ID)

	return s.buildTokenResponse(user)
}

func (s *AuthService) Login(ctx context.Context, req models.LoginRequest) (*models.AuthResponse, error) {
	user, err := s.users.FindByEmail(ctx, req.Email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if user.PasswordHash == nil {
		return nil, ErrInvalidCredentials // OAuth-only account
	}
	if err := bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildTokenResponse(user)
}

func (s *AuthService) Refresh(ctx context.Context, refreshToken string) (*models.AuthResponse, error) {
	claims, err := utils.ParseToken(refreshToken, s.cfg.JWTSecret)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	user, err := s.users.FindByID(ctx, claims.UserID)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	return s.buildTokenResponse(user)
}

func (s *AuthService) Me(ctx context.Context, userID string) (*models.User, error) {
	return s.users.FindByID(ctx, userID)
}

// ForgotPassword generates a new random password, saves it and emails it to the user.
// It never reports whether the email is registered, so the response is the same either way.
func (s *AuthService) ForgotPassword(ctx context.Context, email string) error {
	user, err := s.users.FindByEmail(ctx, email)
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}

	newPassword, err := s.email.GeneratePassword(10)
	if err != nil {
		return err
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	if err := s.users.UpdatePasswordHash(ctx, user.ID, string(hash)); err != nil {
		return err
	}

	name := user.Email
	if user.Name != nil && *user.Name != "" {
		name = *user.Name
	}
	if err := s.email.SendNewPassword(user.Email, name, newPassword); err != nil {
		slog.Error("failed to send forgot-password email", "email", user.Email, "error", err)
		return err
	}
	return nil
}

// SearchUsers does a case-insensitive partial email match, for admin lookup.
func (s *AuthService) SearchUsers(ctx context.Context, emailQuery string) ([]models.User, error) {
	return s.users.SearchByEmail(ctx, emailQuery)
}

// CountUsers returns the total number of registered users, for admin stats.
func (s *AuthService) CountUsers(ctx context.Context) (int, error) {
	return s.users.Count(ctx)
}

func (s *AuthService) buildTokenResponse(user *models.User) (*models.AuthResponse, error) {
	access, refresh, err := utils.GenerateTokenPair(user.ID, user.Email, s.cfg.JWTSecret)
	if err != nil {
		return nil, err
	}
	return &models.AuthResponse{
		AccessToken:  access,
		RefreshToken: refresh,
		User:         user,
	}, nil
}
