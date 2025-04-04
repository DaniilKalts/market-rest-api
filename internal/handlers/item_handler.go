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
	MsgItemDeleted = "item deleted successfully"
)

type ItemHandler struct {
	service services.ItemService
}

func NewItemHandler(service services.ItemService) *ItemHandler {
	return &ItemHandler{service: service}
}

func (h *ItemHandler) HandleCreateItem(ctx *gin.Context) {
	item, err := ginhelpers.GetContextValue[*models.Item](ctx, "model")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.CreateItem(item); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (h *ItemHandler) HandleGetItemByID(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	item, err := h.service.GetItemByID(id)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (h *ItemHandler) HandleGetAllItems(ctx *gin.Context) {
	items, err := h.service.GetAllItems()
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, items)
}

func (h *ItemHandler) HandleUpdateItem(ctx *gin.Context) {
	updateItemDTO, err := ginhelpers.GetContextValue[*models.UpdateItem](
		ctx, "model",
	)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	updatedItem, err := h.service.UpdateItem(id, updateItemDTO)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, updatedItem)
}

func (h *ItemHandler) HandleDeleteItem(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	if err := h.service.DeleteItem(id); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": MsgItemDeleted})
}
