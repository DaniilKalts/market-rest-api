package jwt

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"time"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Role string

const (
	User Role = "user"
)

type Claims struct {
	jwt.RegisteredClaims
	Role Role
}

func generateTokenID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(b), nil
}

func GenerateJWT(issuer string, subject string, minutes uint) (string, error) {
	secret := config.Config.Server.Secret
	if secret == "" {
		return "", errors.New("SECRET is not set in the environment")
	}

	tokenID, err := generateTokenID()
	if err != nil {
		return "", err
	}

	issuedAt := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Duration(minutes) * time.Minute)),
			ID:        tokenID,
		},
		Role: User,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Config.Server.Secret), nil
	})
	if err != nil {
		return claims, err
	}

	return claims, nil
}
