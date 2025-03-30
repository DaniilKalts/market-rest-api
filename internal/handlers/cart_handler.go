package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/pkg/ginhelpers"
	"github.com/DaniilKalts/market-rest-api/pkg/jwt"
)

type CartHandler struct {
	service services.CartService
}

func NewCartHandler(service services.CartService) *CartHandler {
	return &CartHandler{service: service}
}

func (h *CartHandler) HandleAddItem(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(
			http.StatusUnauthorized, gin.H{"error": "invalid user id in token"},
		)
		ctx.Abort()
		return
	}

	cart, err := h.service.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	itemIDStr := ctx.Param("id")

	itemID, err := strconv.Atoi(itemIDStr)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	if err := h.service.AddItem(cart.ID, itemID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "added item"})
}

func (h *CartHandler) HandleUpdateItem(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(
			http.StatusUnauthorized, gin.H{"error": "invalid user id in token"},
		)
		ctx.Abort()
		return
	}

	cart, err := h.service.GetCartByUserID(userID)
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

	item, err := ginhelpers.GetContextValue[*models.UpdateItem](ctx, "model")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdateItem(cart.ID, itemID, item.Quantity); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "updated item"})
}

func (h *CartHandler) HandleDeleteItem(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(
			http.StatusUnauthorized, gin.H{"error": "invalid user id in token"},
		)
		ctx.Abort()
		return
	}

	cart, err := h.service.GetCartByUserID(userID)
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

	if err := h.service.DeleteItem(cart.ID, itemID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "deleted item"})
}

func (h *CartHandler) HandleClearCart(ctx *gin.Context) {
	claims, err := ginhelpers.GetContextValue[*jwt.Claims](ctx, "claims")
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, convErr := strconv.Atoi(claims.Subject)
	if convErr != nil {
		ctx.JSON(
			http.StatusUnauthorized, gin.H{"error": "invalid user id in token"},
		)
		ctx.Abort()
		return
	}

	cart, err := h.service.GetCartByUserID(userID)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if cart == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "cart not found"})
		return
	}

	if err := h.service.ClearCart(cart.ID); err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "cleared cart"})
}
