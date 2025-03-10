package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type CustomClaims struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	jwt.RegisteredClaims
}

func CreateToken(firstName, lastName string) (string, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return "", errors.New("SECRET is not set in the environment")
	}

	expirationTime := time.Now().Add(15 * time.Minute)

	claims := CustomClaims{
		FirstName: firstName,
		LastName:  lastName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return errors.New("SECRET is not set in the environment")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return err
	}

	if !token.Valid {
		return errors.New("Invalid token")
	}

	return nil
}

func ExtractClaimsFromToken(tokenString string) (*CustomClaims, error) {
	secret := os.Getenv("SECRET")
	if secret == "" {
		return nil, errors.New("SECRET is not set in the environment")
	}

	claims := &CustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Unix() < time.Now().Unix() {
		return nil, errors.New("token is expired")
	}
	if claims.FirstName == "" || claims.LastName == "" {
		return nil, errors.New("missing essential claim data")
	}

	return claims, nil
}
