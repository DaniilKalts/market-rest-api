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
	MsgCartCleared = "cart cleared successfully"
)

func getCart(ctx *gin.Context, cartService services.CartService) (
	*models.Cart, error,
) {
	userID, err := getUserIDFromContext(ctx)
	if err != nil {
		return nil, err
	}
	cart, err := cartService.GetCartByUserID(userID)
	if err != nil {
		return nil, err
	}
	if cart == nil {
		return nil, errs.ErrCartNotFound
	}
	return cart, nil
}

type CartHandler struct {
	itemService services.ItemService
	cartService services.CartService
}

func NewCartHandler(
	itemService services.ItemService, cartService services.CartService,
) *CartHandler {
	return &CartHandler{
		itemService: itemService,
		cartService: cartService,
	}
}

func (h *CartHandler) HandleGetCart(ctx *gin.Context) {
	cart, err := getCart(ctx, h.cartService)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cart)
}

func (h *CartHandler) HandleAddItem(ctx *gin.Context) {
	cart, err := getCart(ctx, h.cartService)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(
			http.StatusNotFound, gin.H{"error": errs.ErrItemNotFound.Error()},
		)
		return
	}

	cartItem, err := h.cartService.AddItem(cart.ID, itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, cartItem)
}

func (h *CartHandler) HandleUpdateItem(ctx *gin.Context) {
	cart, err := getCart(ctx, h.cartService)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(
			http.StatusNotFound, gin.H{"error": errs.ErrItemNotFound.Error()},
		)
		return
	}

	updateItem, err := ginhelpers.GetContextValue[*models.UpdateCartItem](
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
	cart, err := getCart(ctx, h.cartService)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	itemIDStr := ctx.Param("id")
	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(
			http.StatusBadRequest, gin.H{"error": errs.ErrInvalidID.Error()},
		)
		return
	}

	item, err := h.itemService.GetItemByID(itemID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if item == nil {
		ctx.JSON(
			http.StatusNotFound, gin.H{"error": errs.ErrItemNotFound.Error()},
		)
		return
	}

	if err := h.cartService.DeleteItem(cart.ID, itemID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": MsgItemDeleted})
}

func (h *CartHandler) HandleClearCart(ctx *gin.Context) {
	cart, err := getCart(ctx, h.cartService)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if err := h.cartService.ClearCart(cart.ID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": MsgCartCleared})
}
