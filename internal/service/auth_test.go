package service

import (
	"AvitoTech/internal/repository"
	mocks "AvitoTech/test/mock"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"AvitoTech/internal/entity"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthService_Authenticate_NewUser(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockToken := new(mocks.MockToken)

	authService := NewAuthService(logger, mockUserRepo, mockToken)

	username := "newuser"
	password := "password"

	mockUserRepo.On("FindUserByUsername", username).Return(&entity.User{}, repository.ErrorUserNotFound)

	newUser := &entity.User{
		Id:       1,
		Username: username,
		Password: "hashedpassword",
		Balance:  1000,
	}
	mockUserRepo.On("InsertUser", mock.AnythingOfType("*entity.User")).Return(newUser, nil)

	mockToken.On("GenerateToken", newUser.Id).Return("generated-token", nil)

	token, err := authService.Authenticate(username, password)

	assert.NoError(t, err)
	assert.Equal(t, "generated-token", token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_Authenticate_ExistingUser(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockToken := new(mocks.MockToken)

	authService := NewAuthService(logger, mockUserRepo, mockToken)

	username := "existinguser"
	password := "validpassword123"

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	existingUser := &entity.User{
		Id:       1,
		Username: username,
		Password: string(hashedPassword),
		Balance:  1000,
	}
	mockUserRepo.On("FindUserByUsername", username).Return(existingUser, nil)

	mockToken.On("GenerateToken", existingUser.Id).Return("generated-token", nil)

	token, err := authService.Authenticate(username, password)

	assert.NoError(t, err)
	assert.Equal(t, "generated-token", token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockToken := new(mocks.MockToken)

	authService := NewAuthService(logger, mockUserRepo, mockToken)

	username := "existinguser"
	password := "wrongpassword"

	existingUser := &entity.User{
		Id:       1,
		Username: username,
		Password: "hashedpassword",
		Balance:  1000,
	}
	mockUserRepo.On("FindUserByUsername", username).Return(existingUser, nil)

	token, err := authService.Authenticate(username, password)

	assert.Error(t, err)
	assert.Empty(t, token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_VerifyJWT_ValidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockToken := new(mocks.MockToken)

	authService := NewAuthService(logger, mockUserRepo, mockToken)

	token := "valid-token"
	userID := 1

	mockToken.On("VerifyToken", token).Return(userID, nil)

	verifiedUserID, err := authService.VerifyJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, userID, verifiedUserID)

	mockToken.AssertExpectations(t)
}

func TestAuthService_VerifyJWT_InvalidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mocks.MockUserRepository)
	mockToken := new(mocks.MockToken)

	authService := NewAuthService(logger, mockUserRepo, mockToken)

	token := "invalid-token"

	mockToken.On("VerifyToken", token).Return(-1, errors.New("invalid token"))

	verifiedUserID, err := authService.VerifyJWT(token)

	assert.Error(t, err)
	assert.Equal(t, -1, verifiedUserID)

	mockToken.AssertExpectations(t)
}
