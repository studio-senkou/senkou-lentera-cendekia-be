package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtManager struct {
	secret []byte
}

type AuthToken struct {
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}

func NewJwtManager(secret string) *JwtManager {
	return &JwtManager{
		secret: []byte(secret),
	}
}

func (j *JwtManager) GenerateToken(userID int, expiry time.Time) (*AuthToken, error) {
	jwtClaims := jwt.MapClaims{
		"payload": userID,
		"exp":     jwt.NewNumericDate(expiry),
		"iat":     jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	signedToken, err := token.SignedString(j.secret)
	if err != nil {
		return nil, fmt.Errorf("failed to sign token: %w", err)
	}

	expTime := jwtClaims["exp"].(*jwt.NumericDate).Time

	return &AuthToken{
		Token:     signedToken,
		ExpiresAt: expTime,
	}, nil
}

func (j *JwtManager) ValidateToken(token string) (jwt.MapClaims, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !parsedToken.Valid {
		return nil, errors.New("token is not valid")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
