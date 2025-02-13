package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	mocks "AvitoTech/test/mock"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestCoinService_SendCoin_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	fromUserID := 1
	toUsername := "receiver"
	amount := 100

	sender := &entity.User{
		Id:       fromUserID,
		Username: "sender",
		Balance:  1000,
	}
	receiver := &entity.User{
		Id:       2,
		Username: toUsername,
		Balance:  500,
	}

	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(receiver, nil)
	mockUserRepo.On("TransferMoney", fromUserID, receiver.Id, amount).Return(nil)
	mockHistoryRepo.On("InsertOperation", entity.Operation{
		FromUser: sender.Username,
		ToUser:   receiver.Username,
		Amount:   amount,
	}).Return(&entity.Operation{Id: 1, FromUser: sender.Username, ToUser: receiver.Username, Amount: amount}, nil)

	err := coinService.SendCoin(fromUserID, toUsername, amount)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_SenderNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	fromUserID := 1
	toUsername := "receiver"
	mockUserRepo.On("FindUserById", fromUserID).Return(&entity.User{}, repository.ErrorUserNotFound)

	err := coinService.SendCoin(fromUserID, toUsername, 100)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrorUserNotFound))

	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_ReceiverNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	fromUserID := 1
	toUsername := "receiver"

	sender := &entity.User{
		Id:       fromUserID,
		Username: "sender",
		Balance:  1000,
	}

	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(&entity.User{}, repository.ErrorUserNotFound)

	err := coinService.SendCoin(fromUserID, toUsername, 100)

	assert.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrorUserNotFound))

	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_TransferFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	fromUserID := 1
	toUsername := "receiver"
	amount := 100

	sender := &entity.User{
		Id:       fromUserID,
		Username: "sender",
		Balance:  1000,
	}
	receiver := &entity.User{
		Id:       2,
		Username: toUsername,
		Balance:  500,
	}

	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(receiver, nil)
	mockUserRepo.On("TransferMoney", fromUserID, receiver.Id, amount).Return(errors.New("transfer failed"))

	err := coinService.SendCoin(fromUserID, toUsername, amount)

	assert.Error(t, err)
	assert.Equal(t, "transfer failed", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	userID := 1
	item := entity.Item{Title: "cup", OwnerId: userID}
	cost := entity.Items[item.Title]

	mockUserRepo.On("WithdrawMoney", userID, cost).Return(nil)
	mockInventoryRepo.On("InsertItem", userID, item.Title).Return(&item, nil)

	err := coinService.BuyItem(userID, item.Title)

	assert.NoError(t, err)

	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_ItemNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	userID := 1
	item := "nonexistent_item"

	err := coinService.BuyItem(userID, item)

	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}

func TestCoinService_BuyItem_WithdrawFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	userID := 1
	item := "cup"
	cost := entity.Items[item]

	mockUserRepo.On("WithdrawMoney", userID, cost).Return(errors.New("insufficient funds"))

	err := coinService.BuyItem(userID, item)

	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_InsertItemFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	userID := 1
	item := entity.Item{Title: "cup", OwnerId: userID}
	cost := entity.Items[item.Title]

	mockUserRepo.On("WithdrawMoney", userID, cost).Return(nil)
	mockInventoryRepo.On("InsertItem", userID, item.Title).Return(nil, errors.New("insert failed"))

	err := coinService.BuyItem(userID, item.Title)

	assert.Error(t, err)
	assert.Equal(t, "insert failed", err.Error())

	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}
