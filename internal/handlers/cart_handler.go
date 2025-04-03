package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
)

type CartHandler struct {
	itemService services.ItemService
	cartService services.CartService
}

func NewCartHandler(
	itemService services.ItemService,
	cartService services.CartService,
) *CartHandler {
	return &CartHandler{itemService: itemService, cartService: cartService}
}

func (h *CartHandler) HandleGetCart(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

func (h *CartHandler) HandleAddItem(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	cartItem, err := h.cartService.AddItem(cart.ID, itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cartItem)
}

func (h *CartHandler) HandleUpdateItem(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	updateItem, err := ginhelpers.GetContextValue[*models.UpdateItem](
		ctx, "model",
	)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	cartItem, err := h.cartService.UpdateItem(
		cart.ID, itemID, updateItem.Quantity,
	)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cartItem)
}

func (h *CartHandler) HandleDeleteItem(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "item not found"})
		return
	}

	if err := h.cartService.DeleteItem(cart.ID, itemID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "item deleted successfully"})
}

func (h *CartHandler) HandleClearCart(ctx *gin.Context) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	cart, err := h.cartService.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	if err := h.cartService.ClearCart(cart.ID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "cart cleared successfully"})
}
