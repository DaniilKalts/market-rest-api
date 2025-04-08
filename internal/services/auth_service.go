package services

import (
	"errors"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"strconv"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

type AuthService interface {
	RegisterUser(user *models.RegisterUser) (string, string, error)
	LoginUser(email, password string) (string, string, error)
	LogoutUser(accessToken, refreshToken string) error
	RefreshTokens(refreshToken string) (string, string, error)
}

type authService struct {
	repo       repositories.UserRepository
	tokenStore redis.TokenStore
}

func NewAuthService(
	repo repositories.UserRepository, tokenStore redis.TokenStore,
) AuthService {
	return &authService{
		repo:       repo,
		tokenStore: tokenStore,
	}
}

func (s *authService) generateAndStoreTokens(userID int, role string) (
	string, string, error,
) {
	uidStr := strconv.Itoa(userID)
	accessToken, err := jwt.GenerateJWT(uidStr, 15, role)
	if err != nil {
		return "", "", errs.ErrTokenGeneration
	}

	refreshToken, err := jwt.GenerateJWT(uidStr, 1440, role)
	if err != nil {
		return "", "", errs.ErrTokenGeneration
	}

	if err := s.tokenStore.SaveJWTokens(
		userID, accessToken, refreshToken,
	); err != nil {
		return "", "", errs.ErrTokenStorage
	}

	return accessToken, refreshToken, nil
}

func (s *authService) RegisterUser(req *models.RegisterUser) (
	string, string, error,
) {
	existingUser, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		if errors.Is(err, errs.ErrUserNotFound) {
			existingUser = nil
		} else {
			return "", "", err
		}
	}
	if existingUser != nil {
		return "", "", errs.ErrUserExists
	}

	user := &models.User{
		Email:       req.Email,
		Password:    req.Password,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		PhoneNumber: req.PhoneNumber,
		Role:        models.RoleUser,
	}

	if err := s.repo.Create(user); err != nil {
		return "", "", errs.ErrUserCreationFailed
	}

	return s.generateAndStoreTokens(user.ID, string(user.Role))
}

func (s *authService) LoginUser(email, password string) (
	string, string, error,
) {
	user, err := s.repo.GetByEmail(email)
	if err != nil {
		return "", "", errs.ErrUserVerifyFailed
	}
	if user == nil {
		return "", "", errs.ErrUserNotFound
	}

	if _, err := jwt.CheckPassword(password, user.Password); err != nil {
		return "", "", errs.ErrInvalidCreds
	}

	return s.generateAndStoreTokens(user.ID, string(user.Role))
}

func (s *authService) LogoutUser(accessToken, refreshToken string) error {
	claims, err := jwt.ParseJWT(accessToken)
	if err != nil {
		return errs.ErrTokenParsingFailed
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return errs.ErrInvalidTokenSub
	}

	if err := s.tokenStore.DeleteJWTokens(
		userID, accessToken, refreshToken,
	); err != nil {
		return errs.ErrTokenDeletionFailed
	}

	return nil
}

func (s *authService) RefreshTokens(refreshToken string) (
	string, string, error,
) {
	claims, err := jwt.ParseJWT(refreshToken)
	if err != nil {
		return "", "", errs.ErrTokenParsingFailed
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return "", "", errs.ErrInvalidTokenSub
	}

	if err := s.tokenStore.DeleteJWToken(userID, refreshToken); err != nil {
		return "", "", errs.ErrTokenDeletionFailed
	}

	accessToken, newRefreshToken, err := s.generateAndStoreTokens(
		userID, claims.Role,
	)
	if err != nil {
		return "", "", err
	}

	if err := s.tokenStore.SaveJWTokens(
		userID, accessToken, newRefreshToken,
	); err != nil {
		return "", "", errs.ErrTokenStorage
	}

	return accessToken, newRefreshToken, nil
}
