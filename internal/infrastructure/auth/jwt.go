package auth

import (
	"time"

	"github.com/contractiq/contractiq/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
)

// Claims represents JWT token claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTService handles JWT token generation and validation.
type JWTService struct {
	secret []byte
	expiry time.Duration
}

// NewJWTService creates a new JWT service.
func NewJWTService(secret string, expiry time.Duration) *JWTService {
	return &JWTService{
		secret: []byte(secret),
		expiry: expiry,
	}
}

// GenerateToken creates a signed JWT for the given user.
func (s *JWTService) GenerateToken(userID, email string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.expiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.secret)
}

// ValidateToken parses and validates a JWT, returning the claims.
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, apperror.NewUnauthorized("invalid token signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return nil, apperror.NewUnauthorized("invalid or expired token")
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, apperror.NewUnauthorized("invalid token claims")
	}

	return claims, nil
}
