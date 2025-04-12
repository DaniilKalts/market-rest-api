package services_test

import (
	mocks2 "github.com/DaniilKalts/market-rest-api/internal/mocks"
	"strconv"
	"testing"
	"time"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

var (
	martinUser = &models.User{
		ID:          1,
		FirstName:   "Martin",
		LastName:    "Kalts",
		Email:       "martin@gmail.com",
		Password:    "$2a$10$83WJqaaF2DpG2d7rTrLYhexbNUiG6tEDmJxVCR.0kdiM26wZXl1Qy",
		PhoneNumber: "+77007473472",
		Role:        models.RoleUser,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}
)

func TestRegisterUser_UserExists(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	req := &models.RegisterUser{
		FirstName:       "Martin",
		LastName:        "Kalts",
		Email:           "martin@gmail.com",
		Password:        "12341234",
		ConfirmPassword: "12341234",
		PhoneNumber:     "+77007473472",
	}

	repoMock.
		On("GetByEmail", req.Email).
		Return(martinUser, nil)

	access, refresh, err := svc.RegisterUser(req)
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, errs.ErrUserExists, err)

	repoMock.AssertExpectations(t)
	tokenStoreMock.AssertExpectations(t)
}

func TestRegisterUser_Success(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	req := &models.RegisterUser{
		FirstName:       "New",
		LastName:        "User",
		Email:           "newuser@example.com",
		Password:        "12341234",
		ConfirmPassword: "12341234",
		PhoneNumber:     "+77007473472",
	}

	repoMock.
		On("GetByEmail", req.Email).
		Return(nil, errs.ErrUserNotFound)

	repoMock.
		On("Create", mock.AnythingOfType("*models.User")).
		Run(
			func(args mock.Arguments) {
				u := args.Get(0).(*models.User)
				u.ID = 2
			},
		).
		Return(nil)

	tokenStoreMock.
		On("SaveJWTokens", 2, mock.Anything, mock.Anything).
		Return(nil)

	access, refresh, err := svc.RegisterUser(req)
	require.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	repoMock.AssertExpectations(t)
	tokenStoreMock.AssertExpectations(t)
}

func TestLoginUser_UserNotFound(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	repoMock.
		On("GetByEmail", "nonexistent@example.com").
		Return(nil, errs.ErrUserNotFound)

	access, refresh, err := svc.LoginUser("nonexistent@example.com", "12341234")
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, errs.ErrUserVerifyFailed, err)

	repoMock.AssertExpectations(t)
}

func TestLoginUser_InvalidCreds(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	repoMock.
		On("GetByEmail", martinUser.Email).
		Return(martinUser, nil)

	access, refresh, err := svc.LoginUser(martinUser.Email, "wrongpass")
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, errs.ErrInvalidCreds, err)

	repoMock.AssertExpectations(t)
}

func TestLoginUser_Success(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	repoMock.
		On("GetByEmail", martinUser.Email).
		Return(martinUser, nil)

	tokenStoreMock.
		On("SaveJWTokens", martinUser.ID, mock.Anything, mock.Anything).
		Return(nil)

	access, refresh, err := svc.LoginUser(martinUser.Email, "12341234")
	require.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	repoMock.AssertExpectations(t)
	tokenStoreMock.AssertExpectations(t)
}

func generateValidToken(userID int, role string, minutes uint) string {
	uidStr := strconv.Itoa(userID)
	token, err := jwt.GenerateJWT(uidStr, minutes, role)
	if err != nil {
		panic("failed to generate token in test: " + err.Error())
	}
	return token
}

func TestLogoutUser_ParseError(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	invalidAccessToken := "invalid.token"
	refreshToken := "dummy-refresh-token"

	err := svc.LogoutUser(invalidAccessToken, refreshToken)
	assert.Equal(t, errs.ErrTokenParsingFailed, err)

	tokenStoreMock.AssertExpectations(t)
}

func TestLogoutUser_DeleteError(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	userID := 1
	accessToken := generateValidToken(userID, string(models.RoleUser), 15)
	refreshToken := generateValidToken(userID, string(models.RoleUser), 1440)

	tokenStoreMock.
		On("DeleteJWTokens", userID, accessToken, refreshToken).
		Return(errs.ErrTokenDeletionFailed)

	err := svc.LogoutUser(accessToken, refreshToken)
	assert.Equal(t, errs.ErrTokenDeletionFailed, err)

	tokenStoreMock.AssertExpectations(t)
}

func TestLogoutUser_Success(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	userID := 1
	accessToken := generateValidToken(userID, string(models.RoleUser), 15)
	refreshToken := generateValidToken(userID, string(models.RoleUser), 1440)

	tokenStoreMock.
		On("DeleteJWTokens", userID, accessToken, refreshToken).
		Return(nil)

	err := svc.LogoutUser(accessToken, refreshToken)
	require.NoError(t, err)

	tokenStoreMock.AssertExpectations(t)
}

func TestRefreshTokens_ParseError(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	invalidRefreshToken := "invalid.token"
	access, refresh, err := svc.RefreshTokens(invalidRefreshToken)
	assert.Empty(t, access)
	assert.Empty(t, refresh)
	assert.Equal(t, errs.ErrTokenParsingFailed, err)
}

func TestRefreshTokens_DeleteError(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	userID := 1
	refreshToken := generateValidToken(userID, string(models.RoleUser), 1440)

	tokenStoreMock.
		On("DeleteJWToken", userID, refreshToken).
		Return(errs.ErrTokenDeletionFailed)

	access, newRefresh, err := svc.RefreshTokens(refreshToken)
	assert.Empty(t, access)
	assert.Empty(t, newRefresh)
	assert.Equal(t, errs.ErrTokenDeletionFailed, err)

	tokenStoreMock.AssertExpectations(t)
}

func TestRefreshTokens_SaveError(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	userID := 1
	refreshToken := generateValidToken(userID, string(models.RoleUser), 1440)

	tokenStoreMock.
		On("DeleteJWToken", userID, refreshToken).
		Return(nil)

	tokenStoreMock.
		On("SaveJWTokens", userID, mock.Anything, mock.Anything).
		Return(errs.ErrTokenStorage)

	access, newRefresh, err := svc.RefreshTokens(refreshToken)
	assert.Empty(t, access)
	assert.Empty(t, newRefresh)
	assert.Equal(t, errs.ErrTokenStorage, err)

	tokenStoreMock.AssertExpectations(t)
}

func TestRefreshTokens_Success(t *testing.T) {
	repoMock := new(mocks2.UserRepository)
	tokenStoreMock := new(mocks2.TokenStore)
	svc := services.NewAuthService(repoMock, tokenStoreMock)

	userID := 1
	oldRefreshToken := generateValidToken(userID, string(models.RoleUser), 1440)

	tokenStoreMock.
		On("DeleteJWToken", userID, oldRefreshToken).
		Return(nil)

	tokenStoreMock.
		On("SaveJWTokens", userID, mock.Anything, mock.Anything).
		Return(nil)

	access, newRefresh, err := svc.RefreshTokens(oldRefreshToken)
	require.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, newRefresh)

	tokenStoreMock.AssertExpectations(t)
}
