package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/barca-strategos/phoenix/pkg/types"
)

// Claims represents JWT claims.
type Claims struct {
	UserID string      `json:"user_id"`
	Role   types.UserRole `json:"role"`
	jwt.RegisteredClaims
}

// Service handles JWT and RBAC.
type Service struct {
	secretKey []byte
}

// New creates an Auth service.
func New(secret string) *Service {
	return &Service{secretKey: []byte(secret)}
}

// GenerateToken creates a JWT for a user.
func (s *Service) GenerateToken(userID string, role types.UserRole) (string, error) {
	claims := Claims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secretKey)
}

// ValidateToken parses and validates a JWT.
func (s *Service) ValidateToken(tokenString string) (*Claims, error) {
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secretKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// RequireRole checks if the user has the required role.
func (s *Service) RequireRole(ctx context.Context, required types.UserRole) error {
	claims, ok := ctx.Value("claims").(*Claims)
	if !ok {
		return fmt.Errorf("unauthenticated")
	}
	if claims.Role != required && claims.Role != types.RoleAdmin {
		return fmt.Errorf("insufficient role")
	}
	return nil
}
