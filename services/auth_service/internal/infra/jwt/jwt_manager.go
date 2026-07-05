package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type jwtConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func newJWTConfig() (*jwtConfig, error) {
	cfg := &jwtConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *jwtConfig) load() {
	c.Secret = env.GetString("JWT_SECRET", "")
	c.AccessTTL = env.GetDuration("JWT_ACCESS_TTL", 15*time.Minute)
	c.RefreshTTL = env.GetDuration("JWT_REFRESH_TTL", 168*time.Hour)
}

func (c *jwtConfig) validate() error {
	if c.Secret == "" {
		return fmt.Errorf("JWT_SECRET %w", errors.ErrIsRequired)
	}
	return nil
}

type jwtManager struct {
	cfg    *jwtConfig
	logger *zap.SugaredLogger
}

func NewManager(logger *zap.SugaredLogger) (domain.TokenManager, error) {
	cfg, err := newJWTConfig()
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	return &jwtManager{
		cfg:    cfg,
		logger: logger,
	}, nil
}

func (m *jwtManager) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	token, err := m.generateJWT(userID, role, m.cfg.AccessTTL)
	if err != nil {
		m.logger.Errorw("failed to generate access token", "error", err)
		return "", fmt.Errorf("generate access token: %w", err)
	}
	return token, nil
}

func (m *jwtManager) GenerateRefreshToken(userID uuid.UUID, role string) (string, error) {
	token, err := m.generateJWT(userID, role, m.cfg.RefreshTTL)
	if err != nil {
		m.logger.Errorw("failed to generate refresh token", "error", err)
		return "", fmt.Errorf("generate refresh token: %w", err)
	}
	return token, nil
}

func (m *jwtManager) Parse(token string) (*domain.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidSigningMethod
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		m.logger.Errorw("failed to parse token", "error", err)
		return nil, domain.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		m.logger.Errorw("failed to get token claims")
		return nil, domain.ErrInvalidToken
	}

	sub, err := m.getStringClaim(claims, "sub")
	if err != nil {
		return nil, err
	}
	role, err := m.getStringClaim(claims, "role")
	if err != nil {
		return nil, err
	}

	userID, err := uuid.Parse(sub)
	if err != nil {
		return nil, err
	}

	return &domain.Claims{
		UserID: userID,
		Role:   role,
	}, nil
}

func (m *jwtManager) GetRefreshTTL() time.Duration {
	return m.cfg.RefreshTTL
}

func (m *jwtManager) generateJWT(userID uuid.UUID, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.cfg.Secret))
	if err != nil {
		return "", fmt.Errorf("signing token: %w", err)
	}

	return signedToken, nil
}

func (m *jwtManager) getStringClaim(claims jwt.MapClaims, key string) (string, error) {
	val, ok := claims[key].(string)
	if !ok {
		m.logger.Errorw("failed to get token claim")
		return "", domain.ErrInvalidClaim
	}
	return val, nil
}
