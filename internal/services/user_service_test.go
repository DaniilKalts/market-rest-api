package services

import (
	stdErrors "errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/mocks"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

func ptr(s string) *string {
	return &s
}

var (
	martinUser = &models.User{
		ID:          1,
		FirstName:   "Martin",
		LastName:    "Kalts",
		Email:       "martin@gmail.com",
		Password:    "$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02",
		PhoneNumber: "+77007473472",
		Role:        models.RoleUser,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	daniilUser = &models.User{
		ID:          2,
		FirstName:   "Daniil",
		LastName:    "Kalts",
		Email:       "daniil.kalts@example.com",
		Password:    "$2a$10$S0meRealHashedPassword",
		PhoneNumber: "+77001234568",
		Role:        models.RoleUser,
		CreatedAt:   time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
	}
)

func TestGetUserByEmail_Success(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetByEmail", martinUser.Email).Return(martinUser, nil)

	userService := NewUserService(mockRepo)

	user, err := userService.GetUserByEmail(martinUser.Email)

	require.NoError(t, err)
	assert.Equal(t, martinUser, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByEmail_RepoError(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	email := "dummy@example.com"

	mockRepo.On("GetByEmail", email).Return(nil, errs.ErrUserNotFound)

	userService := NewUserService(mockRepo)

	user, err := userService.GetUserByEmail(email)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, errs.ErrUserNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_Success(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetByID", martinUser.ID).Return(martinUser, nil)

	userService := NewUserService(mockRepo)

	user, err := userService.GetUserByID(martinUser.ID)

	require.NoError(t, err)
	assert.Equal(t, martinUser, user)

	mockRepo.AssertExpectations(t)
}

func TestGetUserByID_RepoError(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	id := 42

	mockRepo.On("GetByID", id).Return(nil, errs.ErrUserNotFound)

	userService := NewUserService(mockRepo)

	user, err := userService.GetUserByID(id)

	require.Error(t, err)
	assert.Nil(t, user)
	assert.EqualError(t, err, errs.ErrUserNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_NoUsers(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetAll").Return([]models.User{}, nil)

	userService := NewUserService(mockRepo)

	users, err := userService.GetAllUsers()

	require.NoError(t, err)
	assert.Empty(t, users)

	mockRepo.AssertExpectations(t)
}

func TestGetAllUsers_SomeUsers(t *testing.T) {

	expectedUsers := []models.User{
		*martinUser,
		*daniilUser,
	}

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetAll").Return(expectedUsers, nil)

	userService := NewUserService(mockRepo)

	users, err := userService.GetAllUsers()

	require.NoError(t, err)
	assert.Equal(t, expectedUsers, users)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_Success_WithPasswordChange(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	updateDTO := &models.UpdateUser{
		FirstName:       ptr("Martin"),
		LastName:        ptr("Kalts"),
		Email:           ptr("martin.programmer@gmail.com"),
		PhoneNumber:     ptr("+77007473472"),
		Password:        ptr("12341234"),
		ConfirmPassword: ptr("12341234"),
	}

	existingUser := &models.User{
		ID:          1,
		FirstName:   "Martin",
		LastName:    "Kalts",
		Email:       "martin@gmail.com",
		Password:    "$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02",
		PhoneNumber: "+77007473472",
		Role:        models.RoleUser,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	mockRepo.On("GetByID", 1).Return(existingUser, nil)

	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(
		func(u *models.User) *models.User { return u },
		nil,
	)

	userService := NewUserService(mockRepo)

	updatedUser, err := userService.UpdateUserByID(1, updateDTO)

	require.NoError(t, err)
	assert.Equal(t, "Martin", updatedUser.FirstName)
	assert.Equal(t, "Kalts", updatedUser.LastName)
	assert.Equal(t, "martin.programmer@gmail.com", updatedUser.Email)
	assert.Equal(t, "+77007473472", updatedUser.PhoneNumber)
	assert.NotEqual(t, "12341234", updatedUser.Password)

	ok, err := jwt.CheckPassword("12341234", updatedUser.Password)
	require.NoError(t, err)
	assert.True(t, ok)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_Success_WithoutPasswordChange(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	updateDTO := &models.UpdateUser{
		FirstName:   ptr("Martin"),
		LastName:    ptr("Kalts"),
		Email:       ptr("martin.programmer@gmail.com"),
		PhoneNumber: ptr("+77007473472"),
	}

	existingUser := &models.User{
		ID:          1,
		FirstName:   "Martin",
		LastName:    "Kalts",
		Email:       "martin@gmail.com",
		Password:    "$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02",
		PhoneNumber: "+77007473472",
		Role:        models.RoleUser,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	mockRepo.On("GetByID", 1).Return(existingUser, nil)

	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(
		func(u *models.User) *models.User { return u },
		nil,
	)

	userService := NewUserService(mockRepo)

	updatedUser, err := userService.UpdateUserByID(1, updateDTO)

	require.NoError(t, err)
	assert.Equal(t, "Martin", updatedUser.FirstName)
	assert.Equal(t, "Kalts", updatedUser.LastName)
	assert.Equal(t, "martin.programmer@gmail.com", updatedUser.Email)
	assert.Equal(t, "+77007473472", updatedUser.PhoneNumber)
	assert.Equal(
		t, "$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02",
		updatedUser.Password,
	)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_PasswordMismatch(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	existingUser := &models.User{
		ID:          1,
		FirstName:   "Martin",
		LastName:    "Kalts",
		Email:       "martin@gmail.com",
		Password:    "$2a$10$EKq8Yv9Y1WnrDFEdiMYCSOaz/oq2I9l9ngJyH/eBRM3lIbcJRLS02",
		PhoneNumber: "+77007473472",
		Role:        models.RoleUser,
	}

	mockRepo.On("GetByID", 1).Return(existingUser, nil)

	updateDTO := &models.UpdateUser{
		Password:        ptr("12341234"),
		ConfirmPassword: ptr("43214321"),
	}

	userService := NewUserService(mockRepo)

	updatedUser, err := userService.UpdateUserByID(1, updateDTO)

	require.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.EqualError(t, err, "passwords do not match")

	mockRepo.AssertExpectations(t)
}

func TestUpdateUserByID_GetByIDError(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("GetByID", 1).Return(nil, errs.ErrUserNotFound)

	updateDTO := &models.UpdateUser{
		FirstName: ptr("Martin"),
	}

	userService := NewUserService(mockRepo)

	updatedUser, err := userService.UpdateUserByID(1, updateDTO)

	require.Error(t, err)
	assert.Nil(t, updatedUser)
	assert.EqualError(t, err, errs.ErrUserNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestDeleteUserByID_Success(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	mockRepo.On("Delete", 1).Return(nil)

	userService := NewUserService(mockRepo)

	err := userService.DeleteUserByID(1)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUserByID_RepoError(t *testing.T) {

	mockRepo := new(mocks.UserRepository)

	expectedErr := stdErrors.New("delete error")

	mockRepo.On("Delete", 1).Return(expectedErr)

	userService := NewUserService(mockRepo)

	err := userService.DeleteUserByID(1)

	require.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}
