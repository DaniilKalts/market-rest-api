package services

import (
	"errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/repositories"
)

type ItemService interface {
	CreateItem(item *models.Item) error
	GetItemByID(id int) (*models.Item, error)
	GetAllItems() ([]models.Item, error)
	UpdateItem(id int, item *models.UpdateItem) (*models.Item, error)
	DeleteItem(id int) error
}

type itemService struct {
	repo repositories.ItemRepository
}

func NewItemService(repo repositories.ItemRepository) ItemService {
	return &itemService{repo: repo}
}

func (s *itemService) CreateItem(item *models.Item) error {
	return s.repo.Create(item)
}

func (s *itemService) GetItemByID(id int) (*models.Item, error) {
	return s.repo.GetByID(id)
}

func (s *itemService) GetAllItems() ([]models.Item, error) {
	return s.repo.GetAll()
}

func (s *itemService) UpdateItem(
	id int,
	updateItemDTO *models.UpdateItem,
) (*models.Item, error) {
	existingItem, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if existingItem == nil {
		return nil, errors.New("item not found")
	}

	if updateItemDTO.Name != nil {
		existingItem.Name = *updateItemDTO.Name
	}
	if updateItemDTO.Description != nil {
		existingItem.Description = *updateItemDTO.Description
	}
	if updateItemDTO.Price != nil {
		existingItem.Price = *updateItemDTO.Price
	}
	if updateItemDTO.Stock != nil {
		existingItem.Stock = *updateItemDTO.Stock
	}

	err = s.repo.Update(existingItem)
	if err != nil {
		return nil, err
	}
	return existingItem, nil
}

func (s *itemService) DeleteItem(id int) error {
	return s.repo.Delete(id)
}
