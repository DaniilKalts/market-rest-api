package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

var (
	ErrClaimsNotFound = errors.New("claims not found in context")
	ErrInvalidClaims  = errors.New("claims are not of the expected type")
	ErrInvalidUserID  = errors.New("invalid user id in token")
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
	claimsVal, exists := ctx.Get("claims")
	if !exists {
		return 0, ErrClaimsNotFound
	}
	claims, ok := claimsVal.(*jwt.Claims)
	if !ok {
		return 0, ErrInvalidClaims
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, ErrInvalidUserID
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
