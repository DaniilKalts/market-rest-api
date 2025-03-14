package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	userInterface, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.service.CreateUser(user); err != nil {
		logger.Error("CreateUser: Failed to create user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	logger.Info("CreateUser: User created successfully, ID=" + strconv.Itoa(user.ID))
	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("GetUserByID: Invalid User ID: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		logger.Error("GetUserByID: User not found, ID=" + idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.Info("GetUserByID: User retrieved successfully, ID=" + idStr)
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		logger.Error("GetAllUsers: Failed to retrieve users: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve Users"})
		return
	}

	logger.Info("GetAllUsers: Successfully retrieved  " + strconv.Itoa(len(users)) + " users")
	c.JSON(http.StatusOK, users)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	userInterface, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	user, ok := userInterface.(*models.User)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.service.UpdateUser(user); err != nil {
		logger.Error("UpdateUser: Failed to update user " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	logger.Info("UpdateUser: User with ID=" + strconv.Itoa(user.ID) + " updated successfully")
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("DeleteUser: Invalid User ID: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid User ID"})
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		logger.Error("DeleteUser: User not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	logger.Info("DeleteUser: User with ID=" + idStr + " deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
