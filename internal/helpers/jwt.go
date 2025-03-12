package helpers

import (
	"errors"
	"time"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func CreateToken(issuer, subject string) (string, error) {
	secret := config.Config.Server.Secret
	if secret == "" {
		return "", errors.New("SECRET is not set in the environment")
	}

	issuedAt := time.Now()
	expiresAt := issuedAt.Add(15 * time.Minute)

	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
