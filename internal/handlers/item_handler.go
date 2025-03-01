package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/logger"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

type ItemHandler struct {
	service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

// @Summary Create a new item
// @Description Create a new item with the given payload
// @Tags Items
// @Accept json
// @Produce json
// @Param item body models.RequestCreateItem true "Item to create"
// @Success 201 {object} models.Item "Item created"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /item/create [post]
func (h *ItemHandler) CreateItem(c *gin.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("CreateItem: Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := h.service.CreateItem(&item); err != nil {
		logger.Error("CreateItem: Failed to create item: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create item"})
		return
	}

	logger.Info("CreateItem: Item created successfully, ID=" + strconv.Itoa(item.ID))
	c.JSON(http.StatusCreated, item)
}

// @Summary Get item by id
// @Description Get an item with the specified ID
// @Tags Items
// @Accept json
// @Produce json
// @Param id query int true "Item ID"
// @Success 200 {object} models.Item "Item retrieved successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /item [get]
func (h *ItemHandler) GetItemByID(c *gin.Context) {
	idStr := c.Query("id")

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

// @Summary Get all items
// @Description Retrieve a list of all items
// @Tags Items
// @Produce json
// @Success 200 {array} models.Item "A list of items retrieved successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /items [get]
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

// @Summary Update item
// @Description Update an existing item with the given payload
// @Tags Items
// @Accept json
// @Produce json
// @Param item body models.RequestUpdateItem true "Item to update"
// @Success 200 {object} models.Item "Item udpated successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /item/update [put]
func (h *ItemHandler) UpdateItem(c *gin.Context) {
	var item models.Item

	if err := c.ShouldBindJSON(&item); err != nil {
		logger.Error("UpdateItem: Invalid request payload: " + err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if err := h.service.UpdateItem(&item); err != nil {
		logger.Error("UpdateItem: Failed to update item: " + err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update item"})
		return
	}

	logger.Info("UpdateItem: Item with ID=" + strconv.Itoa(item.ID) + " updated successfully")
	c.JSON(http.StatusOK, item)
}

// @Summary Delete item
// @Description Delete an item with the specified ID
// @Tags Items
// @Accept json
// @Produce json
// @Param id query int true "Item ID"
// @Success 200 {object} models.Item "Item deleted successfully"
// @Failure 400 {object} models.BadRequestError "Bad Request"
// @Failure 500 {object} models.InternalServerError "Internal Server Error"
// @Router /item/delete [delete]
func (h *ItemHandler) DeleteItem(c *gin.Context) {
	idStr := c.Query("id")

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
