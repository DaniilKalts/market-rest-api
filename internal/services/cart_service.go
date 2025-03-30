package services

import (
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
)

type CartService interface {
	AddItem(cartID int, itemID int) error
	GetCartByUserID(cartID int) (*models.Cart, error)
	UpdateItem(cartID int, itemID int, quantity uint) error
	DeleteItem(cartID int, itemID int) error
	ClearCart(cartID int) error
}

type cartService struct {
	repo repositories.CartRepository
}

func NewCartService(repo repositories.CartRepository) CartService {
	return &cartService{repo: repo}
}

func (s *cartService) AddItem(cartID int, itemID int) error {
	return s.repo.Add(cartID, itemID)
}

func (s *cartService) GetCartByUserID(userID int) (*models.Cart, error) {
	return s.repo.GetByUserID(userID)
}

func (s *cartService) UpdateItem(cartID int, itemID int, quantity uint) error {
	return s.repo.Update(cartID, itemID, quantity)
}

func (s *cartService) DeleteItem(cartID int, itemID int) error {
	return s.repo.Delete(cartID, itemID)
}

func (s *cartService) ClearCart(cartID int) error {
	return s.repo.Clear(cartID)
}
