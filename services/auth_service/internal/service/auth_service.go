package service

import (
	"context"
	"fmt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	c "github.com/mykytakuzminov/ridely-svc/shared/context"
	log "github.com/mykytakuzminov/ridely-svc/shared/logging"

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
	traceID := c.GetTraceID(ctx)

	// Check existence
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		log.Declined(s.logger, traceID, "registration", err)
		return nil, domain.ErrUserAlreadyExists
	}

	// Hash password
	hpwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		log.Failed(s.logger, traceID, "registration", err)
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
	log.Success(s.logger, traceID, "registration")
	return tokens, nil
}

func (s *authService) Login(ctx context.Context, input *domain.LoginInput) (*domain.TokenPair, error) {
	traceID := c.GetTraceID(ctx)

	// Check existence
	user, err := s.userRepo.GetByEmail(ctx, input.Email)
	if err != nil {
		log.Failed(s.logger, traceID, "login", fmt.Errorf("invalid email"))
		return nil, domain.ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password)); err != nil {
		log.Failed(s.logger, traceID, "login", fmt.Errorf("invalid password"))
		return nil, domain.ErrInvalidCredentials
	}

	// Generate and save tokens
	tokens, err := s.generateAndSaveTokens(ctx, user)
	if err != nil {
		return nil, err
	}

	// Return tokens
	log.Success(s.logger, traceID, "login")
	return tokens, nil
}

func (s *authService) ValidateToken(ctx context.Context, token string) (*domain.Claims, error) {
	claims, err := s.jwtManager.Parse(ctx, token)
	if err != nil {
		return nil, err
	}
	return claims, nil
}

func (s *authService) generateAndSaveTokens(ctx context.Context, user *domain.User) (*domain.TokenPair, error) {
	// Generate tokens
	accessToken, err := s.jwtManager.GenerateAccessToken(ctx, user.ID, user.Role)
	if err != nil {
		return nil, err
	}
	refreshToken, err := s.jwtManager.GenerateRefreshToken(ctx, user.ID, user.Role)
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
