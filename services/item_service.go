package services

import (
	"github.com/DaniilKalts/market-rest-api/models"
	"github.com/DaniilKalts/market-rest-api/repositories"
)

type ItemService interface {
	CreateItem(item *models.Item) error
	GetItemByID(id int) (*models.Item, error)
	GetAllItems() ([]models.Item, error)
	UpdateItem(item *models.Item) error
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

func (s *itemService) UpdateItem(item *models.Item) error {
	return s.repo.Update(item)
}

func (s *itemService) DeleteItem(id int) error {
	return s.repo.Delete(id)
}
