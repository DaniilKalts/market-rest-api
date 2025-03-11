package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/auth"
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
	Email           string `json:"email"`
	Password        string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Register: invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if req.Password != req.ConfirmPassword {
		logger.Error("Register: passwords do not match")
		c.JSON(http.StatusBadRequest, gin.H{"error": "passwords do not match"})
		return
	}

	user := models.User{
		Email:        req.Email,
		PasswordHash: req.Password,
		FirstName:    req.FirstName,
		LastName:     req.LastName,
	}

	if err := h.service.CreateUser(&user); err != nil {
		logger.Error("Register: failed to create user: " + err.Error())
		return
	}

	tokenString, err := auth.CreateToken(user.FirstName, user.LastName)
	if err != nil {
		logger.Info("Register: failed to create token: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	logger.Info("Register: User registered successfully, ID=" + strconv.Itoa(user.ID))
	c.JSON(http.StatusCreated, gin.H{
		"user":         user,
		"access_token": tokenString,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Login: invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	tokenString, err := auth.CreateToken("Daniil", "Kalts")
	if err != nil {
		logger.Error("Login: failed to create token: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenString})
}
