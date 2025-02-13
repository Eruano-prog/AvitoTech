package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"AvitoTech/internal/service"
	mocks "AvitoTech/test/mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInfoService_GetInfo_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)

	infoService := service.NewInfoService(logger, mockUserRepo, mockHistoryRepo, mockInventoryRepo)

	userID := 1
	username := "testuser"
	balance := 1000
	sentOperations := []entity.Operation{
		{Id: 1, FromUser: username, ToUser: "user2", Amount: 100},
	}
	receivedOperations := []entity.Operation{
		{Id: 2, FromUser: "user2", ToUser: username, Amount: 200},
	}
	inventory := map[string]int{
		"item1": 1,
		"item2": 2,
	}

	mockUserRepo.On("FindUserById", userID).Return(&entity.User{
		Id:       userID,
		Username: username,
		Balance:  balance,
	}, nil)
	mockHistoryRepo.On("GetSentByUser", username).Return(sentOperations, nil)
	mockHistoryRepo.On("GetReceivedByUser", username).Return(receivedOperations, nil)
	mockInventoryRepo.On("GetUsersInventory", userID).Return(inventory, nil)

	info, err := infoService.GetInfo(userID)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, balance, info.Coins)
	assert.Equal(t, sentOperations, info.Sent)
	assert.Equal(t, receivedOperations, info.Received)
	assert.Equal(t, inventory, info.Inventory)

	mockUserRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestInfoService_GetInfo_UserNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)

	infoService := service.NewInfoService(logger, mockUserRepo, mockHistoryRepo, mockInventoryRepo)

	userID := 1
	mockUserRepo.On("FindUserById", userID).Return(&entity.User{}, repository.ErrorUserNotFound)

	info, err := infoService.GetInfo(userID)

	assert.Error(t, err)
	assert.Nil(t, info)

	mockUserRepo.AssertExpectations(t)
}

func TestInfoService_GetInfo_HistoryError(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)

	infoService := service.NewInfoService(logger, mockUserRepo, mockHistoryRepo, mockInventoryRepo)

	userID := 1
	username := "testuser"
	balance := 1000
	historyError := errors.New("history error")

	mockUserRepo.On("FindUserById", userID).Return(&entity.User{
		Id:       userID,
		Username: username,
		Balance:  balance,
	}, nil)
	mockHistoryRepo.On("GetSentByUser", username).Return([]entity.Operation{}, historyError)
	mockHistoryRepo.On("GetReceivedByUser", username).Return([]entity.Operation{}, historyError)
	mockInventoryRepo.On("GetUsersInventory", userID).Return(map[string]int{}, nil)

	info, err := infoService.GetInfo(userID)

	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, balance, info.Coins)
	assert.Empty(t, info.Sent)
	assert.Empty(t, info.Received)

	mockUserRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}
