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
	exists, err := s.userRepo.ExistsByEmail(ctx, input.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		s.logger.Warnw("registration rejected, user already exists", "email", input.Email)
		return nil, domain.ErrUserAlreadyExists
	}

	hpwd, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		s.logger.Errorw("registration rejected, failed to hash password", "error", err)
		return nil, err
	}

	user := &domain.User{Email: input.Email, PasswordHash: string(hpwd)}
	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens

	s.logger.Infow("user registered successfully", "user_id", user.ID)
	return nil, nil
}
