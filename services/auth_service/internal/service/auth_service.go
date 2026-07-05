package service

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
)

type authService struct {
	userRepo   domain.UserRepository
	tokenRepo  domain.TokenRepository
	jwtManager domain.TokenManager
	logger     *zap.SugaredLogger
}

func NewAuthService(
	userRepo domain.UserRepository,
	tokenRepo domain.TokenRepository,
	jwtManager domain.TokenManager,
	logger *zap.SugaredLogger,
) domain.AuthService {
	return &authService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		jwtManager: jwtManager,
		logger:     logger.With("layer", "service"),
	}
}

func (s *authService) Register(ctx context.Context, input *domain.RegisterInput) (*domain.TokenPair, error) {
	// Check existence
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		s.logger.Warnw("registration rejected, user already exists", "email", input.Email)
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hpwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		s.logger.Errorw("registration rejected, failed to hash password", "error", err)
		return nil, err
	}

	// Create user
	user := &domain.User{Email: input.Email, PasswordHash: string(hpwd)}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate and save tokens
	tokens, err := s.generateAndSaveTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return tokens
	s.logger.Infow("user registered successfully", "user_id", user.ID)
	return tokens, nil
}

func (s *authService) Login(ctx context.Context, input *domain.LoginInput) (*domain.TokenPair, error) {
	// Check existence
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		s.logger.Warnw("login failed, invalid email", "email", input.Email)
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		s.logger.Warnw("login failed, invalid password")
		return nil, domain.ErrInvalidCredentials
	}

	// Generate and save tokens
	tokens, err := s.generateAndSaveTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return tokens
	s.logger.Infow("user logged in successfully", "user_id", user.ID)
	return tokens, nil
}

func (s *authService) generateAndSaveTokens(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	// Save refresh token
	if err := s.tokenRepo.Set(ctx, &domain.RefreshToken{
		Token:  refreshToken,
		UserID: user.ID,
		TTL:    s.jwtManager.GetRefreshTTL(),
	}); err != nil {
		return nil, err
	}

	// Return token pair
	return &domain.TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
