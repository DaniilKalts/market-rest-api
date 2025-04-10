package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(
	authService services.AuthService,
) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) HandleRegister(ctx *gin.Context) {
	req, err := ginhelpers.GetContextValue[*models.RegisterUser](ctx, "model")
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RegisterUser(req)
	if err != nil {
		_ = ctx.Error(err)
		if errors.Is(err, errs.ErrUserExists) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		}
		return
	}

	if err := jwt.SetAuthCookies(
		ctx.Writer, accessToken, refreshToken,
	); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusCreated, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

func (h *AuthHandler) HandleLogin(ctx *gin.Context) {
	req, err := ginhelpers.GetContextValue[*models.LoginUser](ctx, "model")
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.LoginUser(
		req.Email, req.Password,
	)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.SetAuthCookies(
		ctx.Writer, accessToken, refreshToken,
	); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusOK, gin.H{
			"access_token":  accessToken,
			"refresh_token": refreshToken,
		},
	)
}

func (h *AuthHandler) HandleLogout(ctx *gin.Context) {
	accessToken, err := ctx.Cookie("access_token")
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.LogoutUser(accessToken, refreshToken); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.DeleteAuthCookies(ctx.Writer); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout successfull"})
}

func (h *AuthHandler) HandleRefreshToken(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, newRefreshToken, err := h.service.RefreshTokens(refreshToken)
	if err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.SetAuthCookies(
		ctx.Writer, accessToken, newRefreshToken,
	); err != nil {
		_ = ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(
		http.StatusCreated, gin.H{
			"access_token":  accessToken,
			"refresh_token": newRefreshToken,
		},
	)
}
