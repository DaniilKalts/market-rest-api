package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/logger"
	"github.com/DaniilKalts/market-rest-api/models"
	"github.com/DaniilKalts/market-rest-api/services"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error("CreateUser: Invalid request payload: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateUser(&user); err != nil {
		logger.Error("CreateUser: Failed to create user: " + err.Error())
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
	logger.Info("CreateUser: User created successfully, ID=" + strconv.Itoa(user.ID))
}

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("GetUserByID: Invalid request payload: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		logger.Error("GetUserByID: User not found, ID=" + idStr)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(user)
	logger.Info("GetUserById: User retrieved successfully, ID=" + idStr)
}

func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := h.service.GetAllUsers()
	if err != nil {
		logger.Error("GetAllUsers: Failed to retrieve users: " + err.Error())
		http.Error(w, "Failed to get Users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(users); err != nil {
		logger.Error("GetAllUsers: Failed to encode users: " + err.Error())
		http.Error(w, "Failed to encode users", http.StatusInternalServerError)
		return
	}

	logger.Info("GetAllUsers: Successfully retrieived  " + strconv.Itoa(len(users)) + " users")
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		logger.Error("UpdateUser: Invalid request payload: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateUser(&user); err != nil {
		logger.Error("UpdateUser: Failed to update user " + err.Error())
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
	logger.Info("UpdateUser: User with ID=" + strconv.Itoa(user.ID) + " was updated successfully")
}

func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("DeleteUser: Invalid User ID: " + err.Error())
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteUser(id); err != nil {
		logger.Error("DeleteUser: User not found: " + err.Error())
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("DeleteUser: User with ID=" + idStr + " was deleted successfully")
}
