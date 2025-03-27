package services

import (
	"errors"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

type AuthService interface {
	RegisterAndAuthenticateUser(user *models.RegisterUser) (string, string, error)
	AuthenticateUser(email string, password string) (string, string, error)
	LogoutUser(accessToken string, refreshToken string) error
	RefreshToken(refreshToken string) (string, string, error)
}

type authService struct {
	repo       repositories.UserRepository
	tokenStore *redis.TokenStore
}

func NewAuthService(repo repositories.UserRepository, tokenStore *redis.TokenStore) AuthService {
	return &authService{repo: repo, tokenStore: tokenStore}
}

func (r *authService) generateAndSaveTokens(userID int, role string) (string, string, error) {
	accessToken, err := jwt.GenerateJWT(strconv.Itoa(userID), 15, role)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := jwt.GenerateJWT(strconv.Itoa(userID), 1440, role)
	if err != nil {
		return "", "", err
	}

	if err := r.tokenStore.SaveJWTokens(userID, accessToken, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func (r *authService) RegisterAndAuthenticateUser(req *models.RegisterUser) (string, string, error) {
	existingUser, err := r.repo.GetByEmail(req.Email)
	if err != nil {
		return "", "", err
	}
	if existingUser != nil {
		return "", "", errors.New("user already exists")
	}

	user := &models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Role:      "user",
	}

	if err := r.repo.Create(user); err != nil {
		return "", "", err
	}

	return r.generateAndSaveTokens(user.ID, string(user.Role))
}

func (r *authService) AuthenticateUser(email string, password string) (string, string, error) {
	user, err := r.repo.GetByEmail(email)
	if err != nil {
		return "", "", err
	}
	if user == nil {
		return "", "", errors.New("user not found")
	}

	if _, err := jwt.CheckPassword(password, user.Password); err != nil {
		return "", "", err
	}

	return r.generateAndSaveTokens(user.ID, string(user.Role))
}

func (r *authService) LogoutUser(accessToken string, refreshToken string) error {
	claims, err := jwt.ParseJWT(accessToken)
	if err != nil {
		return err
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return convErr
	}

	if err := r.tokenStore.DeleteJWTokens(userID, accessToken, refreshToken); err != nil {
		return err
	}

	return nil
}

func (r *authService) RefreshToken(refreshToken string) (string, string, error) {
	claims, err := jwt.ParseJWT(refreshToken)
	if err != nil {
		return "", "", err
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		return "", "", convErr
	}

	if err := r.tokenStore.DeleteJWToken(userID, refreshToken); err != nil {
		return "", "", err
	}

	accessToken, refreshToken, err := r.generateAndSaveTokens(userID, claims.Role)
	if err != nil {
		return "", "", err
	}

	if err := r.tokenStore.SaveJWTokens(userID, accessToken, refreshToken); err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
