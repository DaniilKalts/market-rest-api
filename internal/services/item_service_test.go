package services

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"
	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/mocks"
)

func TestCreateItem_Success(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	item := &models.Item{
		ID:          1,
		Name:        "T-shirt",
		Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
		Price:       30,
		Stock:       20,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	mockRepo.On("Create", item).Return(nil)

	itemService := NewItemService(mockRepo)

	err := itemService.CreateItem(item)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestCreateItem_Error(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	item := &models.Item{
		ID:          1,
		Name:        "T-shirt",
		Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
		Price:       30,
		Stock:       20,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	expectedErr := errors.New("create error")

	mockRepo.On("Create", item).Return(expectedErr)

	itemService := NewItemService(mockRepo)

	err := itemService.CreateItem(item)

	require.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetItemByID_Success(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	item := &models.Item{
		ID:          1,
		Name:        "T-shirt",
		Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
		Price:       30,
		Stock:       20,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	mockRepo.On("GetByID", item.ID).Return(item, nil)

	itemService := NewItemService(mockRepo)

	result, err := itemService.GetItemByID(item.ID)

	require.NoError(t, err)
	assert.Equal(t, item, result)

	mockRepo.AssertExpectations(t)
}

func TestGetItemByID_RepoError(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	id := 42

	mockRepo.On("GetByID", id).Return(nil, errs.ErrItemNotFound)

	itemService := NewItemService(mockRepo)

	result, err := itemService.GetItemByID(id)

	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, errs.ErrItemNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestGetAllItems_Success(t *testing.T) {

	expectedItems := []models.Item{
		{
			ID:          1,
			Name:        "T-shirt",
			Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
			Price:       30,
			Stock:       20,
			CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
			UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		},
		{
			ID:          2,
			Name:        "Sweater",
			Description: "A comfortable and stylish sweater made from high-quality materials.",
			Price:       50,
			Stock:       10,
			CreatedAt:   time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
			UpdatedAt:   time.Date(2025, 3, 1, 10, 0, 0, 0, time.UTC),
		},
	}

	mockRepo := new(mocks.ItemRepository)

	mockRepo.On("GetAll").Return(expectedItems, nil)

	itemService := NewItemService(mockRepo)

	items, err := itemService.GetAllItems()

	require.NoError(t, err)
	assert.Equal(t, expectedItems, items)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_Success(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	existingItem := &models.Item{
		ID:          1,
		Name:        "T-shirt",
		Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
		Price:       30,
		Stock:       20,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	updateDTO := &models.UpdateItem{
		Name:        ptr("T-shirt Updated"),
		Description: ptr("A premium quality T-shirt featuring an exclusive IITU logo design."),
		Price:       ptrUint(35),
		Stock:       ptrUint(15),
	}

	mockRepo.On("GetByID", existingItem.ID).Return(existingItem, nil)

	mockRepo.On("Update", mock.AnythingOfType("*models.Item")).Return(nil)

	itemService := NewItemService(mockRepo)

	updatedItem, err := itemService.UpdateItem(existingItem.ID, updateDTO)

	require.NoError(t, err)
	assert.Equal(t, "T-shirt Updated", updatedItem.Name)
	assert.Equal(
		t, "A premium quality T-shirt featuring an exclusive IITU logo design.",
		updatedItem.Description,
	)
	assert.Equal(t, uint(35), updatedItem.Price)
	assert.Equal(t, uint(15), updatedItem.Stock)

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_ItemNotFound(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	id := 42

	mockRepo.On("GetByID", id).Return(nil, errs.ErrItemNotFound)

	updateDTO := &models.UpdateItem{
		Name: ptr("T-shirt Updated"),
	}

	itemService := NewItemService(mockRepo)

	updatedItem, err := itemService.UpdateItem(id, updateDTO)

	require.Error(t, err)
	assert.Nil(t, updatedItem)
	assert.EqualError(t, err, errs.ErrItemNotFound.Error())

	mockRepo.AssertExpectations(t)
}

func TestUpdateItem_UpdateError(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	existingItem := &models.Item{
		ID:          1,
		Name:        "T-shirt",
		Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
		Price:       30,
		Stock:       20,
		CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
		UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	}

	updateDTO := &models.UpdateItem{
		Name: ptr("T-shirt Updated"),
	}

	mockRepo.On("GetByID", existingItem.ID).Return(existingItem, nil)

	expectedErr := errors.New("update error")
	mockRepo.On(
		"Update", mock.AnythingOfType("*models.Item"),
	).Return(expectedErr)

	itemService := NewItemService(mockRepo)

	updatedItem, err := itemService.UpdateItem(existingItem.ID, updateDTO)

	require.Error(t, err)
	assert.Nil(t, updatedItem)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestDeleteItem_Success(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	mockRepo.On("Delete", 1).Return(nil)

	itemService := NewItemService(mockRepo)

	err := itemService.DeleteItem(1)

	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestDeleteItem_RepoError(t *testing.T) {

	mockRepo := new(mocks.ItemRepository)

	expectedErr := errors.New("delete error")
	mockRepo.On("Delete", 1).Return(expectedErr)

	itemService := NewItemService(mockRepo)

	err := itemService.DeleteItem(1)

	require.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func ptrUint(u uint) *uint {
	return &u
}
