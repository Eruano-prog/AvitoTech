package repository

import (
	"AvitoTech/internal/entity"
	"github.com/stretchr/testify/mock"
)

// MockUserRepository is a mock implementation of the UserRepository interface
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

func (m *MockUserRepository) FindUserById(id int) (*entity.User, error) {
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
