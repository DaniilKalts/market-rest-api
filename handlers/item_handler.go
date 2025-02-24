package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/DaniilKalts/market-rest-api/logger"
	"github.com/DaniilKalts/market-rest-api/models"
	"github.com/DaniilKalts/market-rest-api/services"
)

type ItemHandler struct {
	service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) CreateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item

	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Error("CreateItem: Invalid request payload: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateItem(&item); err != nil {
		logger.Error("CreateItem: Failed to create item: " + err.Error())
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(item)
	logger.Info("CreateItem: Item created successfully, ID=" + strconv.Itoa(item.ID))
}

func (h *ItemHandler) GetItemByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("GetItemByID: Invalid Item ID: " + err.Error())
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	item, err := h.service.GetItemByID(id)
	if err != nil {
		logger.Error("GetItemByID: Item not found, ID=" + idStr)
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(item)
	logger.Info("GetItemByID: Item retrieved successfully, ID=" + idStr)
}

func (h *ItemHandler) GetAllItems(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAllItems()
	if err != nil {
		logger.Error("GetAllItems: Failed to retrieve items: " + err.Error())
		http.Error(w, "Failed to retrieve items", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(items); err != nil {
		logger.Error("GetAllItems: Failed to encode items: " + err.Error())
		http.Error(w, "Failed to encode items", http.StatusInternalServerError)
		return
	}

	logger.Info("GetAllItems: Successfully retrieved " + strconv.Itoa(len(items)) + " items")
}

func (h *ItemHandler) UpdateItem(w http.ResponseWriter, r *http.Request) {
	var item models.Item
	if err := json.NewDecoder(r.Body).Decode(&item); err != nil {
		logger.Error("UpdateItem: Invalid request payload: " + err.Error())
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.UpdateItem(&item); err != nil {
		logger.Error("UpdateItem: Failed to update item: " + err.Error())
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(item)
	logger.Info("UpdateItem: Item with ID=" + strconv.Itoa(item.ID) + " was updated successfully")
}

func (h *ItemHandler) DeleteItem(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("DeleteItem: Invalid Item ID: " + err.Error())
		http.Error(w, "Invalid Item ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteItem(id); err != nil {
		logger.Error("DeleteItem: Item not found: " + err.Error())
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	logger.Info("DeleteItem: Item with ID=" + idStr + " was deleted successfully")
}
