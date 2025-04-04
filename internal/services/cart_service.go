package services

import (
	"errors"
	"fmt"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
)

type CartService interface {
	AddItem(cartID int, itemID int) (*models.CartItem, error)
	GetCartByUserID(cartID int) (*models.Cart, error)
	UpdateItem(cartID int, itemID int, quantity uint) (*models.CartItem, error)
	DeleteItem(cartID int, itemID int) error
	ClearCart(cartID int) error
}

type cartService struct {
	repo        repositories.CartRepository
	itemService ItemService
}

func NewCartService(
	repo repositories.CartRepository,
	itemService ItemService,
) CartService {
	return &cartService{
		repo:        repo,
		itemService: itemService,
	}
}

func (s *cartService) AddItem(cartID int, itemID int) (
	*models.CartItem,
	error,
) {
	item, err := s.itemService.GetItemByID(itemID)
	if err != nil {
		return nil, err
	}
	if item == nil {
		return nil, errors.New("item not found")
	}

	existingCartItem, err := s.repo.GetCartItem(cartID, itemID)
	currentQuantity := uint(0)
	if err == nil && existingCartItem != nil {
		currentQuantity = existingCartItem.Quantity
	}

	if currentQuantity+1 > item.Stock {
		return nil, fmt.Errorf(
			"available stock is %d and you already have %d in your cart",
			item.Stock,
			currentQuantity,
		)
	}

	return s.repo.Add(cartID, itemID)
}

func (s *cartService) GetCartByUserID(userID int) (
	*models.Cart, error,
) {
	return s.repo.GetByUserID(userID)
}

func (s *cartService) UpdateItem(
	cartID int,
	itemID int,
	quantity uint,
) (*models.CartItem, error) {
	return s.repo.Update(cartID, itemID, quantity)
}

func (s *cartService) DeleteItem(cartID int, itemID int) error {
	return s.repo.Delete(cartID, itemID)
}

func (s *cartService) ClearCart(cartID int) error {
	return s.repo.Clear(cartID)
}
