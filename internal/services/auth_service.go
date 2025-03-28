package services

import (
	"errors"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

var (
	// User related errors
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrUserNotFound           = errors.New("user not found")
	ErrUserCreationFailed     = errors.New("failed to create user")
	ErrUserVerificationFailed = errors.New("failed to verify user")
	ErrInvalidCredentials     = errors.New("invalid credentials")

	// Token related errors
	ErrTokenGeneration     = errors.New("failed to generate token")
	ErrTokenStorage        = errors.New("failed to store token")
	ErrTokenParsingFailed  = errors.New("failed to validate token")
	ErrInvalidTokenSubject = errors.New("invalid token subject")
	ErrTokenDeletionFailed = errors.New("failed to delete token")
	ErrTokenSaveFailed     = errors.New("failed to save token")
)

type AuthService interface {
	RegisterUser(user *models.RegisterUser) (string, string, error)
	LoginUser(email string, password string) (string, string, error)
	LogoutUser(accessToken string, refreshToken string) error
	RefreshTokens(refreshToken string) (string, string, error)
}

type authService struct {
	repo       repositories.UserRepository
	tokenStore *redis.TokenStore
}

func NewAuthService(repo repositories.UserRepository, tokenStore *redis.TokenStore) AuthService {
	return &authService{repo: repo, tokenStore: tokenStore}
}

func (r *authService) generateAndStoreTokens(userID int, role string) (string, string, error) {
	accessToken, err := jwt.GenerateJWT(strconv.Itoa(userID), 15, role)
	if err != nil {
		return "", "", ErrTokenGeneration
	}

	refreshToken, err := jwt.GenerateJWT(strconv.Itoa(userID), 1440, role)
	if err != nil {
		return "", "", ErrTokenGeneration
	}

	if err := r.tokenStore.SaveJWTokens(userID, accessToken, refreshToken); err != nil {
		return "", "", ErrTokenStorage
	}

	return accessToken, refreshToken, nil
}

func (r *authService) RegisterUser(req *models.RegisterUser) (string, string, error) {
	existingUser, err := r.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, repositories.ErrUserNotFound) {
			existingUser = nil
		} else {
			return "", "", err
		}
	}
	if existingUser != nil {
		return "", "", ErrUserAlreadyExists
	}

	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      models.RoleUser,
	}

	if err := r.repo.Create(user); err != nil {
		return "", "", ErrUserCreationFailed
	}

	return r.generateAndStoreTokens(user.ID, string(user.Role))
}

func (r *authService) LoginUser(email string, password string) (string, string, error) {
	user, err := r.repo.GetByEmail(email)
	if err != nil {
		return "", "", ErrUserVerificationFailed
	}
	if user == nil {
		return "", "", ErrUserNotFound
	}

	if _, err := jwt.CheckPassword(password, user.Password); err != nil {
		return "", "", ErrInvalidCredentials
	}

	return r.generateAndStoreTokens(user.ID, string(user.Role))
}

func (r *authService) LogoutUser(accessToken string, refreshToken string) error {
	claims, err := jwt.ParseJWT(accessToken)
	if err != nil {
		return ErrTokenParsingFailed
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return ErrInvalidTokenSubject
	}

	if err := r.tokenStore.DeleteJWTokens(userID, accessToken, refreshToken); err != nil {
		return ErrTokenDeletionFailed
	}

	return nil
}

func (r *authService) RefreshTokens(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseJWT(refreshToken)
	if err != nil {
		return "", "", ErrTokenParsingFailed
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return "", "", ErrInvalidTokenSubject
	}

	if err := r.tokenStore.DeleteJWToken(userID, refreshToken); err != nil {
		return "", "", ErrTokenDeletionFailed
	}

	accessToken, refreshToken, err := r.generateAndStoreTokens(userID, claims.Role)
	if err != nil {
		return "", "", err
	}

	if err := r.tokenStore.SaveJWTokens(userID, accessToken, refreshToken); err != nil {
		return "", "", ErrTokenSaveFailed
	}

	return accessToken, refreshToken, nil
}
