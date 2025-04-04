package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
)

const (
	MsgUserDeleted = "user deleted successfully"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) HandleGetUserByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (h *UserHandler) HandleGetAllUsers(ctx *gin.Context) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, users)
}

func (h *UserHandler) HandleUpdateUserByID(ctx *gin.Context) {
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

	idStr := ctx.Param("id")
	userID, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	updatedUser, err := h.service.UpdateUserByID(userID, updateUser)
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

func (h *UserHandler) HandleDeleteUser(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	if err := h.service.DeleteUserByID(id); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": MsgUserDeleted})
}
