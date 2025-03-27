package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/models"
)

var ErrItemNotFound = errors.New("item not found")

type ItemRepository interface {
	Create(item *models.Item) error
	GetByID(id int) (*models.Item, error)
	GetAll() ([]models.Item, error)
	Update(item *models.Item) error
	Delete(id int) error
}

type itemRepository struct {
	db *gorm.DB
}

func NewItemRepository(db *gorm.DB) ItemRepository {
	return &itemRepository{db: db}
}

func (r *itemRepository) Create(item *models.Item) error {
	return r.db.Create(item).Error
}

func (r *itemRepository) GetByID(id int) (*models.Item, error) {
	var item models.Item

	err := r.db.First(&item, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}

	return &item, nil
}

func (r *itemRepository) GetAll() ([]models.Item, error) {
	var items []models.Item

	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (r *itemRepository) Update(item *models.Item) error {
	return r.db.Save(item).Error
}

func (r *itemRepository) Delete(id int) error {
	result := r.db.Delete(&models.Item{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}
