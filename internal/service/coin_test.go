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

	// Мокируем данные
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

	// Мокируем вызовы
	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(receiver, nil)
	mockUserRepo.On("TransferMoney", fromUserID, receiver.Id, amount).Return(nil)
	mockHistoryRepo.On("InsertOperation", entity.Operation{
		FromUser: sender.Username,
		ToUser:   receiver.Username,
		Amount:   amount,
	}).Return(nil)

	// Вызываем метод
	err := coinService.SendCoin(fromUserID, toUsername, amount)

	// Проверяем результаты
	assert.NoError(t, err)

	// Проверяем, что все моки были вызваны
	mockUserRepo.AssertExpectations(t)
	mockHistoryRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_SenderNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем ошибку
	fromUserID := 1
	toUsername := "receiver"
	mockUserRepo.On("FindUserById", fromUserID).Return(&entity.User{}, repository.ErrorUserNotFound)

	// Вызываем метод
	err := coinService.SendCoin(fromUserID, toUsername, 100)

	// Проверяем результаты
	assert.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrorUserNotFound))

	// Проверяем, что мок был вызван
	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_ReceiverNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
	fromUserID := 1
	toUsername := "receiver"

	sender := &entity.User{
		Id:       fromUserID,
		Username: "sender",
		Balance:  1000,
	}

	// Мокируем ошибку
	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(&entity.User{}, repository.ErrorUserNotFound)

	// Вызываем метод
	err := coinService.SendCoin(fromUserID, toUsername, 100)

	// Проверяем результаты
	assert.Error(t, err)
	assert.True(t, errors.Is(err, repository.ErrorUserNotFound))

	// Проверяем, что моки были вызваны
	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_SendCoin_TransferFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
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

	// Мокируем ошибку
	mockUserRepo.On("FindUserById", fromUserID).Return(sender, nil)
	mockUserRepo.On("FindUserByUsername", toUsername).Return(receiver, nil)
	mockUserRepo.On("TransferMoney", fromUserID, receiver.Id, amount).Return(errors.New("transfer failed"))

	// Вызываем метод
	err := coinService.SendCoin(fromUserID, toUsername, amount)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Equal(t, "transfer failed", err.Error())

	// Проверяем, что моки были вызваны
	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_Success(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
	userID := 1
	item := "cup"
	cost := entity.Items[item] // Предположим, что item1 существует и его стоимость определена

	// Мокируем вызовы
	mockUserRepo.On("WithdrawMoney", userID, cost).Return(nil)
	mockInventoryRepo.On("InsertItem", userID, item).Return(nil)

	// Вызываем метод
	err := coinService.BuyItem(userID, item)

	// Проверяем результаты
	assert.NoError(t, err)

	// Проверяем, что все моки были вызваны
	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_ItemNotFound(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
	userID := 1
	item := "nonexistent_item" // Предположим, что такого предмета нет

	// Вызываем метод
	err := coinService.BuyItem(userID, item)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Equal(t, "item not found", err.Error())
}

func TestCoinService_BuyItem_WithdrawFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
	userID := 1
	item := "cup"
	cost := entity.Items[item] // Предположим, что item1 существует и его стоимость определена

	// Мокируем ошибку
	mockUserRepo.On("WithdrawMoney", userID, cost).Return(errors.New("insufficient funds"))

	// Вызываем метод
	err := coinService.BuyItem(userID, item)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Equal(t, "insufficient funds", err.Error())

	// Проверяем, что мок был вызван
	mockUserRepo.AssertExpectations(t)
}

func TestCoinService_BuyItem_InsertItemFailed(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockInventoryRepo := new(mocks.MockInventoryRepository)
	mockHistoryRepo := new(mocks.MockHistoryRepository)

	coinService := NewCoinService(logger, mockUserRepo, mockInventoryRepo, mockHistoryRepo)

	// Мокируем данные
	userID := 1
	item := "cup"
	cost := entity.Items[item] // Предположим, что item1 существует и его стоимость определена

	// Мокируем вызовы
	mockUserRepo.On("WithdrawMoney", userID, cost).Return(nil)
	mockInventoryRepo.On("InsertItem", userID, item).Return(errors.New("insert failed"))

	// Вызываем метод
	err := coinService.BuyItem(userID, item)

	// Проверяем результаты
	assert.Error(t, err)
	assert.Equal(t, "insert failed", err.Error())

	// Проверяем, что моки были вызваны
	mockUserRepo.AssertExpectations(t)
	mockInventoryRepo.AssertExpectations(t)
}
