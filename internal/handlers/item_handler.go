package handlers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type ItemHandler struct {
	service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) CreateItem(c *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](c, "claims")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if claims.Role != "admin" {
		err := errors.New("admin access only")
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	item, err := ginhelpers.GetContextValue[*models.Item](c, "model")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateItem(item); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) GetItemByID(c *gin.Context) {
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.service.GetItemByID(id)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) GetAllItems(c *gin.Context) {
	items, err := h.service.GetAllItems()
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, items)
}

func (h *ItemHandler) UpdateItem(c *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](c, "claims")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if claims.Role != "admin" {
		err := errors.New("admin access only")
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	item, err := ginhelpers.GetContextValue[*models.Item](c, "model")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateItem(item); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func (h *ItemHandler) DeleteItem(c *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](c, "claims")
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if claims.Role != "admin" {
		err := errors.New("admin access only")
		c.Error(err)
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
		return
	}

	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.service.DeleteItem(id); err != nil {
		c.Error(err)
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}
