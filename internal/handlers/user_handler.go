package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/logger"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// @Summary Create a new user
// @Description Create a new user with the given payload
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.RequestCreateUser true "User to create"
// @Success 201 {object} models.RequestCreateUser "User created"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /user/create [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("CreateUser: Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		logger.Error("CreateUser: Failed to create user: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	logger.Info("CreateUser: User created successfully, ID=" + strconv.Itoa(user.ID))
	c.JSON(http.StatusCreated, user)
}

// @Summary Get user by id
// @Description Get a user with the specified ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} models.User "User retrieved successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /user [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Query("id")

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

// @Summary Get all users
// @Description Retrieve a list of all users
// @Tags Users
// @Produce json
// @Success 200 {array} models.User "A list of users retrieved successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /users [get]
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

// @Summary Update user
// @Description Update an existing user with the given payload
// @Tags Users
// @Accept json
// @Produce json
// @Param user body models.RequestUpdateUser true "User to update"
// @Success 200 {object} models.User "User udpated successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /user/update [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		logger.Error("UpdateUser: Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := h.service.UpdateUser(&user); err != nil {
		logger.Error("UpdateUser: Failed to update user " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
		return
	}

	logger.Info("UpdateUser: User with ID=" + strconv.Itoa(user.ID) + " updated successfully")
	c.JSON(http.StatusOK, user)
}

// @Summary Delete user
// @Description Delete a user with the specified ID
// @Tags Users
// @Accept json
// @Produce json
// @Param id query int true "User ID"
// @Success 200 {object} models.User "User deleted successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /user/delete [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Query("id")

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
