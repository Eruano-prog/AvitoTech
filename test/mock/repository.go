// Package mock
package mock

import (
	"AvitoTech/internal/entity"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) InsertUser(user *entity.User) (*entity.User, error) {
	args := m.Called(user)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByUsername(username string) (*entity.User, error) {
	args := m.Called(username)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) FindUserByID(id int) (*entity.User, error) {
	args := m.Called(id)
	return args.Get(0).(*entity.User), args.Error(1)
}

func (m *MockUserRepository) TransferMoney(userFrom int, userTo int, amount int) error {
	args := m.Called(userFrom, userTo, amount)
	return args.Error(0)
}

func (m *MockUserRepository) WithdrawMoney(user int, amount int) error {
	args := m.Called(user, amount)
	return args.Error(0)
}

type MockHistoryRepository struct {
	mock.Mock
}

func (m *MockHistoryRepository) InsertOperation(operation entity.Operation) (*entity.Operation, error) {
	args := m.Called(operation)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Operation), args.Error(1)
}

func (m *MockHistoryRepository) GetSentByUser(name string) ([]entity.Operation, error) {
	args := m.Called(name)
	return args.Get(0).([]entity.Operation), args.Error(1)
}

func (m *MockHistoryRepository) GetReceivedByUser(name string) ([]entity.Operation, error) {
	args := m.Called(name)
	return args.Get(0).([]entity.Operation), args.Error(1)
}

func (m *MockHistoryRepository) DeleteOperation(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

type MockInventoryRepository struct {
	mock.Mock
}

func (m *MockInventoryRepository) InsertItem(owner int, item string) (*entity.Item, error) {
	args := m.Called(owner, item)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Item), args.Error(1)
}

func (m *MockInventoryRepository) GetUsersInventory(userID int) (map[string]int, error) {
	args := m.Called(userID)
	return args.Get(0).(map[string]int), args.Error(1)
}

func (m *MockInventoryRepository) DeleteItem(id int) error {
	args := m.Called(id)
	return args.Error(0)
}
