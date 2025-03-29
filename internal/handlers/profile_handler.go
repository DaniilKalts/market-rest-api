package handlers

import (
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type ProfileHandler struct {
	userService services.UserService
	authService services.AuthService
}

func NewProfileHandler(userService services.UserService, authService services.AuthService) *ProfileHandler {
	return &ProfileHandler{userService: userService, authService: authService}
}

func (h *ProfileHandler) HandleGetProfile(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		ctx.Abort()
		return
	}

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *ProfileHandler) HandleUpdateProfile(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		ctx.Abort()
		return
	}

	user, err := ginhelpers.GetContextValue[*models.UpdateUser](ctx, "model")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := user.Validate(); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.userService.UpdateUserByID(userID, user); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "profile updated successfully"})
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

	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user id in token"})
		ctx.Abort()
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
