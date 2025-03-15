package auth

import (
	"errors"
	"time"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
}

func GenerateJWT(issuer string, subject string, minutes uint) (string, error) {
	secret := config.Config.Server.Secret
	if secret == "" {
		return "", errors.New("SECRET is not set in the environment")
	}

	issuedAt := time.Now()
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   subject,
			IssuedAt:  jwt.NewNumericDate(issuedAt),
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(time.Duration(minutes) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ParseJWT(tokenString string) (jwt.MapClaims, error) {
	claims := jwt.MapClaims{}

	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(config.Config.Server.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	return claims, nil
}
