package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/logger"
)

type ItemHandler struct {
	service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	itemInterface, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	item, ok := itemInterface.(*models.Item)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.service.CreateItem(item); err != nil {
		logger.Error("CreateItem: Failed to create item: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	logger.Info("CreateItem: Item created successfully, ID=" + strconv.Itoa(item.ID))
	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetItemByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("GetItemByID: Invalid Item ID: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
		return
	}

	item, err := h.service.GetItemByID(id)
	if err != nil {
		logger.Error("GetItemByID: Item not found, ID=" + idStr)
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	logger.Info("GetItemByID: Item retrieved successfully, ID=" + idStr)
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.service.GetAllItems()
	if err != nil {
		logger.Error("GetAllItems: Failed to retrieve items: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve items"})
		return
	}

	logger.Info("GetAllItems: Successfully retrieved " + strconv.Itoa(len(items)) + " items")
	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	itemInterface, exists := c.Get("model")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	item, ok := itemInterface.(*models.Item)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request payload"})
		return
	}

	if err := h.service.UpdateItem(item); err != nil {
		logger.Error("UpdateItem: Failed to update item: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	logger.Info("UpdateItem: Item with ID=" + strconv.Itoa(item.ID) + " updated successfully")
	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		logger.Error("DeleteItem: Invalid Item ID: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Item ID"})
		return
	}
	if err := h.service.DeleteItem(id); err != nil {
		logger.Error("DeleteItem: Item not found: " + err.Error())
		c.JSON(http.StatusNotFound, gin.H{"error": "Item not found"})
		return
	}

	logger.Info("DeleteItem: Item with ID=" + idStr + " deleted successfully")
	c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
}
