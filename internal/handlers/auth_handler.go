package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"github.com/DaniilKalts/market-rest-api/internal/config"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/auth"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(authService services.AuthService) *AuthHandler {
	return &AuthHandler{service: authService}
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	reqInterface, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req, ok := reqInterface.(*RegisterRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.Email == "" || req.Password == "" || req.ConfirmPassword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "email, password and confirm password are required"})
		return
	}
	if req.Password != req.ConfirmPassword {
		logger.Error("Register: passwords do not match")
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.service.RegisterUser(&user); err != nil {
		logger.Error("Register: failed to create user: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	accessToken, refreshToken, err := auth.SetAuthCookies(c.Writer, user.ID)
	if err != nil {
		logger.Error("Register: failed to set auth cookies: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set auth cookies"})
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	req, ok := reqInterface.(*LoginRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := h.service.AuthenticateUser(req.Email, req.Password)
	if err != nil {
		logger.Error("Login: invalid credentials: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, err := auth.SetAuthCookies(c.Writer, user.ID)
	if err != nil {
		logger.Error("Login: failed to set auth cookies: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set auth cookies"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	refreshCookie, err := c.Cookie("refresh_token")
	if err != nil {
		logger.Error("RefreshToken: refresh token missing: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token missing"})
		return
	}

	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(refreshCookie, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Invalid signing method")
		}
		return []byte(config.Config.Server.Secret), nil
	})
	if err != nil {
		logger.Error("RefreshToken: error parsing token: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired refresh token"})
		c.Abort()
		return
	}

	userIDStr, ok := claims["sub"].(string)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID, convErr := strconv.Atoi(userIDStr)
	if convErr != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID in token"})
		return
	}

	accessToken, refreshToken, err := auth.SetAuthCookies(c.Writer, userID)
	if err != nil {
		logger.Error("Register: failed to set auth cookies: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to set auth cookies"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}
