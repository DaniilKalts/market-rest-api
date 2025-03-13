package handlers

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/helpers"
	"github.com/DaniilKalts/market-rest-api/internal/logger"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

type AuthHandler struct {
	service services.UserService
}

func NewAuthHandler(userService services.UserService) *AuthHandler {
	return &AuthHandler{service: userService}
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

	_, err := h.service.GetUserByEmail(req.Email)
	if err == nil {
		logger.Error("Register: user already exists")
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		logger.Error("Register: error checking for user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error checking for user"})
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  req.Password,
		FirstName: req.FirstName,
		LastName:  req.LastName,
	}

	if err := h.service.CreateUser(&user); err != nil {
		logger.Error("Register: failed to create user: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to create user"})
		return
	}

	accessToken, refreshToken, err := helpers.SetAuthCookies(c.Writer, user.ID)
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
	}

	req, ok := reqInterface.(*LoginRequest)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, err := h.service.GetUserByEmail(req.Email)
	if err != nil {
		logger.Error("Login: invalid credentials: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logger.Error("Login: invalid credentials: " + err.Error())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	accessToken, refreshToken, err := helpers.SetAuthCookies(c.Writer, user.ID)
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
