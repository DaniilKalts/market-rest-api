package repositories

import (
	"errors"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"gorm.io/gorm"
)

type CartRepository interface {
	Add(cartID int, itemID int) error
	Update(cartID int, itemID int, quantity uint) error
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
) error {
	var cartItem models.CartItem

	err := r.db.Where("cart_id = ? AND item_id = ?", cartID, itemID).
		First(&cartItem).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return r.db.
			Create(
				&models.CartItem{
					CartID:   cartID,
					ItemID:   itemID,
					Quantity: 1,
				},
			).
			Error
	} else if err != nil {
		return err
	}
	cartItem.Quantity += 1

	return r.db.Save(&cartItem).Error
}

func (r *cartRepository) Update(
	cartID int, itemID int, quantity uint,
) error {
	return r.db.
		Where("cart_id = ? AND item_id = ?", cartID, itemID).
		Update("quantity", quantity).Error
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
