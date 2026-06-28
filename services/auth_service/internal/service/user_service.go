package service

import (
	"context"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
)

type userService struct {
	repo   domain.UserRepository
	logger *zap.SugaredLogger
}

func NewUserService(repo domain.UserRepository, logger *zap.SugaredLogger) domain.UserService {
	return &userService{
		repo:   repo,
		logger: logger.With("layer", "service"),
	}
}

func (s *userService) Register(ctx context.Context, input *domain.RegisterInput) (*domain.TokenPair, error) {
	exists, err := s.repo.ExistsByEmail(ctx, input.Email)
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
	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	// Generate tokens

	s.logger.Infow("user registered successfully", "user_id", user.ID)
	return nil, nil
}
