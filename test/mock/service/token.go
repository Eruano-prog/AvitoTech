package service

import "github.com/stretchr/testify/mock"

// MockToken is a mock implementation of the Token interface
type MockToken struct {
	mock.Mock
}

func (m *MockToken) GenerateToken(userID int) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockToken) VerifyToken(tokenString string) (int, error) {
	args := m.Called(tokenString)
	return args.Int(0), args.Error(1)
}
