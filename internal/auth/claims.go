package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrMissingSecret is returned when JWT signing is attempted without a secret.
	ErrMissingSecret = errors.New("auth: missing jwt secret")
	// ErrMissingUserID is returned when a token would be issued without a user id.
	ErrMissingUserID = errors.New("auth: missing user id")
)

// Claims describes the authenticated user carried in an Orbit JWT.
type Claims struct {
	UserID string   `json:"user_id"`
	Roles  []string `json:"roles"`
	jwt.RegisteredClaims
}

// Issue signs a short-lived HS256 token for a user.
func Issue(secret string, userID string, ttl time.Duration) (string, error) {
	if secret == "" {
		return "", ErrMissingSecret
	}
	if userID == "" {
		return "", ErrMissingUserID
	}
	if ttl <= 0 {
		return "", errors.New("auth: ttl must be positive")
	}

	now := time.Now().UTC()
	claims := Claims{
		UserID: userID,
		Roles:  []string{"user"},
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(ttl)),
		},
	}

	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
}

// Parse validates an HS256 token and returns its claims.
func Parse(secret, token string) (*Claims, error) {
	if secret == "" {
		return nil, ErrMissingSecret
	}
	if token == "" {
		return nil, errors.New("auth: missing token")
	}

	claims := &Claims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("auth: invalid token")
	}
	if claims.UserID == "" {
		claims.UserID = claims.Subject
	}
	if claims.UserID == "" {
		return nil, ErrMissingUserID
	}

	return claims, nil
}
