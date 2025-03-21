package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
	"github.com/DaniilKalts/market-rest-api/pkg/redis"
)

type AuthHandler struct {
	service    services.AuthService
	tokenStore *redis.TokenStore
}

func NewAuthHandler(authService services.AuthService, tokenStore *redis.TokenStore) *AuthHandler {
	return &AuthHandler{service: authService, tokenStore: tokenStore}
}

func (h *AuthHandler) Register(c *gin.Context) {
	reqInterface, exists := c.Get("model")
	if !exists {
		err := errors.New("request payload not found")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req, ok := reqInterface.(*models.RegisterUser)
	if !ok {
		err := errors.New("invalid user payload")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := req.Validate(); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.service.RegisterUser(&user); err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := jwt.SetAuthCookies(c.Writer, user.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenStore.SaveJWTokens(user.ID, accessToken, refreshToken); err != nil {
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
	reqInterface, exists := c.Get("model")
	if !exists {
		err := errors.New("invalid user payload")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req, ok := reqInterface.(*models.LoginUser)
	if !ok {
		err := errors.New("invalid user payload")
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := jwt.SetAuthCookies(c.Writer, user.ID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenStore.SaveJWTokens(user.ID, accessToken, refreshToken); err != nil {
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

	accessTokenClaims, err := jwt.ParseJWT(accessToken)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(accessTokenClaims.Subject)
	if convErr != nil {
		c.Error(convErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": convErr.Error()})
		return
	}

	if err := jwt.DeleteAuthCookies(c.Writer); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenStore.DeleteJWTokens(userID, accessToken, refreshToken); err != nil {
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

	claims, err := jwt.ParseJWT(refreshToken)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		c.Error(convErr)
		c.JSON(http.StatusUnauthorized, gin.H{"error": convErr.Error()})
		return
	}

	accessToken, refreshToken, err := jwt.SetAuthCookies(c.Writer, userID)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.tokenStore.SaveJWTokens(userID, accessToken, refreshToken); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
