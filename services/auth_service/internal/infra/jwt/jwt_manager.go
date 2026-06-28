package jwt

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/mykytakuzminov/ridely-svc/services/auth_service/internal/domain"
	"github.com/mykytakuzminov/ridely-svc/shared/env"
	"github.com/mykytakuzminov/ridely-svc/shared/errors"
)

type JWTConfig struct {
	Secret     string
	AccessTTL  time.Duration
	RefreshTTL time.Duration
}

func NewJWTConfig() (*JWTConfig, error) {
	cfg := &JWTConfig{}
	cfg.load()
	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func (c *JWTConfig) load() {
	c.Secret = env.GetString("JWT_SECRET", "")
	c.AccessTTL = env.GetDuration("JWT_ACCESS_TTL", 15*time.Minute)
	c.RefreshTTL = env.GetDuration("JWT_REFRESH_TTL", 168*time.Hour)
}

func (c *JWTConfig) validate() error {
	if c.Secret == "" {
		return fmt.Errorf("JWT_SECRET %w", errors.ErrIsRequired)
	}
	return nil
}

type JWTManager struct {
	cfg *JWTConfig
}

func NewJWTManager(cfg *JWTConfig) domain.TokenManager {
	return &JWTManager{cfg: cfg}
}

func (m *JWTManager) GenerateAccessToken(userID uuid.UUID, role string) (string, error) {
	token, err := m.generateJWT(userID, role, m.cfg.AccessTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *JWTManager) GenerateRefreshToken(userID uuid.UUID, role string) (string, error) {
	token, err := m.generateJWT(userID, role, m.cfg.RefreshTTL)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (m *JWTManager) Parse(token string) (*domain.Claims, error) {
	parsedToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrInvalidSigningMethod
		}
		return []byte(m.cfg.Secret), nil
	})
	if err != nil {
		return nil, domain.ErrInvalidToken
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
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

func (m *JWTManager) generateJWT(userID uuid.UUID, role string, ttl time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"exp":  time.Now().Add(ttl).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(m.cfg.Secret))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func (m *JWTManager) getStringClaim(claims jwt.MapClaims, key string) (string, error) {
	val, ok := claims[key].(string)
	if !ok {
		return "", domain.ErrInvalidClaim
	}
	return val, nil
}
