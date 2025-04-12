package services_test

import (
	"errors"
	"github.com/DaniilKalts/market-rest-api/internal/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errs "github.com/DaniilKalts/market-rest-api/internal/errors"

	"github.com/DaniilKalts/market-rest-api/internal/models"
	"github.com/DaniilKalts/market-rest-api/internal/services"
)

var now = time.Now()

var sampleItem = &models.Item{
	ID:          1,
	Name:        "T-shirt",
	Description: "A premium quality T-shirt featuring an exclusive IITU logo design, crafted from soft, breathable fabric for both style and everyday comfort.",
	Price:       30,
	Stock:       20,
	CreatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
	UpdatedAt:   time.Date(2025, 2, 25, 12, 37, 32, 0, time.UTC),
}

func TestItem_Create_Success(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	mockRepo.On("Create", sampleItem).Return(nil).Once()

	itemService := services.NewItemService(mockRepo)
	err := itemService.CreateItem(sampleItem)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestItem_Create_Error(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	expectedErr := errors.New("create error")
	mockRepo.On("Create", sampleItem).Return(expectedErr).Once()

	itemService := services.NewItemService(mockRepo)
	err := itemService.CreateItem(sampleItem)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestItem_GetByID_Success(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	mockRepo.On("GetByID", sampleItem.ID).Return(sampleItem, nil).Once()

	itemService := services.NewItemService(mockRepo)
	result, err := itemService.GetItemByID(sampleItem.ID)
	require.NoError(t, err)
	assert.Equal(t, sampleItem, result)

	mockRepo.AssertExpectations(t)
}

func TestItem_GetByID_Error(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	id := 42
	repoErr := errs.ErrItemNotFound
	mockRepo.On("GetByID", id).Return(nil, repoErr).Once()

	itemService := services.NewItemService(mockRepo)
	result, err := itemService.GetItemByID(id)
	require.Error(t, err)
	assert.Nil(t, result)
	assert.EqualError(t, err, repoErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestItem_GetAll_Success(t *testing.T) {
	expectedItems := []models.Item{
		*sampleItem,
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
	mockRepo.On("GetAll").Return(expectedItems, nil).Once()

	itemService := services.NewItemService(mockRepo)
	items, err := itemService.GetAllItems()
	require.NoError(t, err)
	assert.Equal(t, expectedItems, items)

	mockRepo.AssertExpectations(t)
}

func TestItem_Update_Success(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	updateDTO := &models.UpdateItem{
		Name:        ptr("T-shirt Updated"),
		Description: ptr("Updated description."),
		Price:       ptrUint(35),
		Stock:       ptrUint(15),
	}

	mockRepo.On("GetByID", sampleItem.ID).Return(sampleItem, nil).Once()
	mockRepo.On(
		"Update", mock.AnythingOfType("*models.Item"),
	).Return(nil).Once()

	itemService := services.NewItemService(mockRepo)
	updatedItem, err := itemService.UpdateItem(sampleItem.ID, updateDTO)
	require.NoError(t, err)
	assert.Equal(t, "T-shirt Updated", updatedItem.Name)
	assert.Equal(t, "Updated description.", updatedItem.Description)
	assert.Equal(t, uint(35), updatedItem.Price)
	assert.Equal(t, uint(15), updatedItem.Stock)

	mockRepo.AssertExpectations(t)
}

func TestItem_Update_GetByIDError(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	id := sampleItem.ID
	expectedErr := errors.New("get error")
	mockRepo.On("GetByID", id).Return(nil, expectedErr).Once()

	updateDTO := &models.UpdateItem{
		Name: ptr("T-shirt Updated"),
	}

	itemService := services.NewItemService(mockRepo)
	updatedItem, err := itemService.UpdateItem(id, updateDTO)
	require.Error(t, err)
	assert.Nil(t, updatedItem)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestItem_Update_ItemNotFound(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	id := 42
	mockRepo.On("GetByID", id).Return(nil, nil).Once()

	updateDTO := &models.UpdateItem{
		Name: ptr("T-shirt Updated"),
	}

	itemService := services.NewItemService(mockRepo)
	updatedItem, err := itemService.UpdateItem(id, updateDTO)
	require.Error(t, err)
	assert.Nil(t, updatedItem)
	assert.EqualError(t, err, "item not found")

	mockRepo.AssertExpectations(t)
}

func TestItem_Update_UpdateError(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)

	updateDTO := &models.UpdateItem{Name: ptr("T-shirt Updated")}
	mockRepo.On("GetByID", sampleItem.ID).Return(sampleItem, nil).Once()

	expectedErr := errors.New("update error")
	mockRepo.On(
		"Update", mock.AnythingOfType("*models.Item"),
	).Return(expectedErr).Once()

	itemService := services.NewItemService(mockRepo)
	updatedItem, err := itemService.UpdateItem(sampleItem.ID, updateDTO)
	require.Error(t, err)
	assert.Nil(t, updatedItem)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func TestItem_Delete_Success(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)
	mockRepo.On("Delete", sampleItem.ID).Return(nil).Once()

	itemService := services.NewItemService(mockRepo)
	err := itemService.DeleteItem(sampleItem.ID)
	require.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestItem_Delete_Error(t *testing.T) {
	mockRepo := new(mocks.ItemRepository)
	expectedErr := errors.New("delete error")
	mockRepo.On("Delete", sampleItem.ID).Return(expectedErr).Once()

	itemService := services.NewItemService(mockRepo)
	err := itemService.DeleteItem(sampleItem.ID)
	require.Error(t, err)
	assert.EqualError(t, err, expectedErr.Error())

	mockRepo.AssertExpectations(t)
}

func ptr(s string) *string {
	return &s
}

func ptrUint(u uint) *uint {
	return &u
}
