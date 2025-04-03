package repositories

import (
	"errors"

	"gorm.io/gorm"

	"github.com/DaniilKalts/market-rest-api/internal/models"
)

var ErrCartNotFound = errors.New("cart not found")

type CartRepository interface {
	Add(cartID int, itemID int) (*models.CartItem, error)
	GetByUserID(userID int) (*models.Cart, error)
	Update(cartID int, itemID int, quantity uint) (*models.CartItem, error)
	Delete(cartID int, itemID int) error
	Clear(cartID int) error
}

type cartRepository struct {
	db *gorm.DB
}

func NewCartRepository(db *gorm.DB) CartRepository {
	return &cartRepository{db: db}
}

func (r *cartRepository) Add(
	cartID int, itemID int,
) (*models.CartItem, error) {
	var cartItem models.CartItem

	err := r.db.
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		cartItem = models.CartItem{
			CartID:   cartID,
			ItemID:   itemID,
			Quantity: 1,
		}
		if err := r.db.Create(&cartItem).Error; err != nil {
			return nil, err
		}
		return &cartItem, nil
	} else if err != nil {
		return nil, err
	}

	cartItem.Quantity++
	if err := r.db.Save(&cartItem).Error; err != nil {
		return nil, err
	}

	if err := r.db.
		Preload("Item").
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepository) GetByUserID(userID int) (*models.Cart, error) {
	var cart models.Cart

	err := r.db.Where("user_id = ?", userID).
		Preload("Items.Item").
		First(&cart).
		Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, ErrCartNotFound
	}

	return &cart, err
}

func (r *cartRepository) Update(
	cartID int,
	itemID int,
	quantity uint,
) (*models.CartItem, error) {
	var cartItem models.CartItem

	err := r.db.
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrItemNotFound
	} else if err != nil {
		return nil, err
	}

	if err = r.db.
		Model(&cartItem).
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		Update("quantity", quantity).Error; err != nil {
		return nil, err
	}

	if err = r.db.
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}

	if err = r.db.
		Preload("Item").
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error; err != nil {
		return nil, err
	}

	return &cartItem, nil
}

func (r *cartRepository) Delete(
	cartID int, itemID int,
) error {
	return r.db.
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		Delete(&models.CartItem{}).
		Error
}

func (r *cartRepository) Clear(cartID int) error {
	return r.db.
		Where("cart_id = ?", cartID).
		Delete(&models.CartItem{}).
		Error
}
