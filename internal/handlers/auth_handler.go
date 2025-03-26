package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(authService services.AuthService, tokenStore *redis.TokenStore) *AuthHandler {
	return &AuthHandler{service: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	req, err := ginhelpers.GetContextValue[*models.RegisterUser](c, "model")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RegisterAndAuthenticateUser(req)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.SetAuthCookies(c.Writer, accessToken, refreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	req, err := ginhelpers.GetContextValue[*models.LoginUser](c, "model")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.SetAuthCookies(c.Writer, accessToken, refreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.LogoutUser(accessToken, refreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.DeleteAuthCookies(c.Writer); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "logout successfull"})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := h.service.RefreshToken(refreshToken)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := jwt.SetAuthCookies(c.Writer, accessToken, refreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
