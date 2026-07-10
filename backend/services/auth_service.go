package services

import (
	"context"
	"errors"

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
	cfg   *config.Config
}

func NewAuthService(
	users repository.UserRepository,
	subs repository.SubscriptionRepository,
	cfg *config.Config,
) *AuthService {
	return &AuthService{users: users, subs: subs, cfg: cfg}
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

// SearchUsers does a case-insensitive partial email match, for admin lookup.
func (s *AuthService) SearchUsers(ctx context.Context, emailQuery string) ([]models.User, error) {
	return s.users.SearchByEmail(ctx, emailQuery)
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
