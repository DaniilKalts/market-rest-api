package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type ProfileHandler struct {
	userService services.UserService
	authService services.AuthService
}

func NewProfileHandler(
	userService services.UserService, authService services.AuthService,
) *ProfileHandler {
	return &ProfileHandler{userService: userService, authService: authService}
}

func getUserIDFromContext(ctx *gin.Context) (int, error) {
	val, exists := ctx.Get("userID")
	if !exists {
		return 0, errors.New("user id not found in context")
	}
	userID, ok := val.(int)
	if !ok {
		return 0, errors.New("user id is not an int")
	}
	return userID, nil
}

func (h *ProfileHandler) HandleGetProfile(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	userResponse := models.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (h *ProfileHandler) HandleUpdateProfile(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	updateUser, err := ginhelpers.GetContextValue[*models.UpdateUser](
		ctx, "model",
	)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := updateUser.Validate(); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updatedUser, err := h.userService.UpdateUserByID(userID, updateUser)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	userResponse := models.UserResponse{
		ID:          updatedUser.ID,
		FirstName:   updatedUser.FirstName,
		LastName:    updatedUser.LastName,
		Email:       updatedUser.Email,
		PhoneNumber: updatedUser.PhoneNumber,
	}

	ctx.JSON(http.StatusOK, userResponse)
}

func (h *ProfileHandler) HandleDeleteProfile(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.authService.LogoutUser(accessToken, refreshToken); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.DeleteUserByID(userID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.DeleteAuthCookies(ctx.Writer); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile deleted successfully"})
}
