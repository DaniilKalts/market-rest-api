package services_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
	"github.com/DaniilKalts/market-rest-api/mocks"
)

type itemServiceStub struct {
	item *models.Item
	err  error
}

func (s *itemServiceStub) CreateItem(item *models.Item) error {
	return nil
}

func (s *itemServiceStub) GetItemByID(id int) (
	*models.Item,
	error,
) {
	return s.item, s.err
}

func (s *itemServiceStub) GetAllItems() (
	[]models.Item,
	error,
) {
	return nil, nil
}

func (s *itemServiceStub) UpdateItem(
	id int,
	updateItemDTO *models.UpdateItem,
) (*models.Item, error) {
	return nil, nil
}

func (s *itemServiceStub) DeleteItem(id int) error {
	return nil
}

var now = time.Now()

var (
	sampleItem = &models.Item{
		ID:    42,
		Name:  "Test Item",
		Stock: 5,
	}
	sampleCartItem = &models.CartItem{
		CartID:    1,
		ItemID:    42,
		Quantity:  2,
		CreatedAt: now,
		UpdatedAt: now,
	}
	sampleCart = &models.Cart{
		ID:        1,
		UserID:    1,
		Items:     []models.CartItem{*sampleCartItem},
		CreatedAt: now,
		UpdatedAt: now,
	}
)

func TestAddItem_Err(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	someErr := fmt.Errorf("service error")
	itemService := &itemServiceStub{item: nil, err: someErr}
	cartService := services.NewCartService(mockRepo, itemService)

	cartItem, err := cartService.AddItem(1, 42)
	assert.Nil(t, cartItem)
	assert.EqualError(t, err, someErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestAddItem_NotFound(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{item: nil, err: nil}
	cartService := services.NewCartService(mockRepo, itemService)

	cartItem, err := cartService.AddItem(1, 42)
	assert.Nil(t, cartItem)
	assert.EqualError(t, err, errs.ErrItemNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestAddItem_Success(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{item: sampleItem, err: nil}
	cartService := services.NewCartService(mockRepo, itemService)

	mockRepo.On("GetCartItem", 1, 42).Return(
		nil, errors.New("not found"),
	).Once()
	mockRepo.On("Add", 1, 42).Return(sampleCartItem, nil).Once()

	cartItem, err := cartService.AddItem(1, 42)
	assert.NoError(t, err)
	assert.Equal(t, sampleCartItem, cartItem)

	mockRepo.AssertExpectations(t)
}

func TestAddItem_ExceedStock(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{
		item: &models.Item{
			ID: 42, Name: "Test Item", Stock: 3,
		},
	}
	cartService := services.NewCartService(mockRepo, itemService)

	existing := &models.CartItem{
		CartID: 1, ItemID: 42, Quantity: 3, CreatedAt: now, UpdatedAt: now,
	}
	mockRepo.On("GetCartItem", 1, 42).Return(existing, nil).Once()

	cartItem, err := cartService.AddItem(1, 42)
	assert.Nil(t, cartItem)
	expectedErrMsg := fmt.Sprintf(
		"available stock is %d and you already have %d in your cart", 3, 3,
	)
	assert.EqualError(t, err, expectedErrMsg)

	mockRepo.AssertExpectations(t)
}

func TestGetCartByUserID_Success(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{}
	cartService := services.NewCartService(mockRepo, itemService)

	mockRepo.On("GetByUserID", 1).Return(sampleCart, nil).Once()

	cart, err := cartService.GetCartByUserID(1)
	require.NoError(t, err)
	assert.Equal(t, sampleCart, cart)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_Success(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{item: sampleItem}
	cartService := services.NewCartService(mockRepo, itemService)

	updated := &models.CartItem{
		CartID:    sampleCartItem.CartID,
		ItemID:    sampleCartItem.ItemID,
		Quantity:  4,
		CreatedAt: sampleCartItem.CreatedAt,
		UpdatedAt: sampleCartItem.UpdatedAt,
	}
	mockRepo.On("Update", 1, 42, uint(4)).Return(updated, nil).Once()

	result, err := cartService.UpdateItem(1, 42, 4)
	require.NoError(t, err)
	assert.Equal(t, updated, result)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_ExceedStock(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{
		item: &models.Item{
			ID: 42, Name: "Test Item", Stock: 5,
		},
	}
	cartService := services.NewCartService(mockRepo, itemService)

	result, err := cartService.UpdateItem(1, 42, 6)
	assert.Nil(t, result)
	expectedErrMsg := fmt.Sprintf(
		"requested quantity %d exceeds available stock %d", 6, 5,
	)
	assert.EqualError(t, err, expectedErrMsg)
	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateItem_Err(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	someErr := fmt.Errorf("service error")
	itemService := &itemServiceStub{item: nil, err: someErr}
	cartService := services.NewCartService(mockRepo, itemService)

	cartItem, err := cartService.UpdateItem(1, 42, 6)
	assert.Nil(t, cartItem)
	assert.EqualError(t, err, someErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_NotFound(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{item: nil, err: nil}
	cartService := services.NewCartService(mockRepo, itemService)

	cartItem, err := cartService.UpdateItem(1, 42, 6)
	assert.Nil(t, cartItem)
	assert.EqualError(t, err, errs.ErrItemNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestDeleteItem_Success(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{}
	cartService := services.NewCartService(mockRepo, itemService)

	mockRepo.On("Delete", 1, 42).Return(nil).Once()
	err := cartService.DeleteItem(1, 42)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestClearCart_Success(t *testing.T) {
	mockRepo := new(mocks.CartRepository)
	itemService := &itemServiceStub{}
	cartService := services.NewCartService(mockRepo, itemService)

	mockRepo.On("Clear", 1).Return(nil).Once()
	err := cartService.ClearCart(1)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}
