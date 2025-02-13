package service

import (
	"AvitoTech/internal/repository"
	mockRepo "AvitoTech/test/mock/repository"
	mockService "AvitoTech/test/mock/service"
	"errors"
	"golang.org/x/crypto/bcrypt"
	"testing"

	"AvitoTech/internal/entity"
	"AvitoTech/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

func TestAuthService_Authenticate_NewUser(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockToken := new(mockService.MockToken)

	authService := service.NewAuthService(logger, mockUserRepo, mockToken)

	username := "newuser"
	password := "password"

	// Mock the FindUserByUsername to return a "user not found" error
	mockUserRepo.On("FindUserByUsername", username).Return(&entity.User{}, repository.ErrorUserNotFound)

	// Mock the InsertUser to return a new user
	newUser := &entity.User{
		Id:       1,
		Username: username,
		Password: "hashedpassword",
		Balance:  1000,
	}
	mockUserRepo.On("InsertUser", mock.AnythingOfType("*entity.User")).Return(newUser, nil)

	// Mock the GenerateToken to return a token
	mockToken.On("GenerateToken", newUser.Id).Return("generated-token", nil)

	token, err := authService.Authenticate(username, password)

	assert.NoError(t, err)
	assert.Equal(t, "generated-token", token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_Authenticate_ExistingUser(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockToken := new(mockService.MockToken)

	authService := service.NewAuthService(logger, mockUserRepo, mockToken)

	username := "existinguser"
	password := "validpassword123" // Используем более длинный пароль

	// Хэшируем пароль для мока
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	assert.NoError(t, err)

	// Mock the FindUserByUsername to return an existing user
	existingUser := &entity.User{
		Id:       1,
		Username: username,
		Password: string(hashedPassword), // Используем хэшированный пароль
		Balance:  1000,
	}
	mockUserRepo.On("FindUserByUsername", username).Return(existingUser, nil)

	// Mock the GenerateToken to return a token
	mockToken.On("GenerateToken", existingUser.Id).Return("generated-token", nil)

	token, err := authService.Authenticate(username, password)

	assert.NoError(t, err)
	assert.Equal(t, "generated-token", token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_Authenticate_InvalidPassword(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockToken := new(mockService.MockToken)

	authService := service.NewAuthService(logger, mockUserRepo, mockToken)

	username := "existinguser"
	password := "wrongpassword"

	// Mock the FindUserByUsername to return an existing user
	existingUser := &entity.User{
		Id:       1,
		Username: username,
		Password: "hashedpassword",
		Balance:  1000,
	}
	mockUserRepo.On("FindUserByUsername", username).Return(existingUser, nil)

	// Mock the bcrypt.CompareHashAndPassword to return an error (invalid password)
	token, err := authService.Authenticate(username, password)

	assert.Error(t, err)
	assert.Empty(t, token)

	mockUserRepo.AssertExpectations(t)
	mockToken.AssertExpectations(t)
}

func TestAuthService_VerifyJWT_ValidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockToken := new(mockService.MockToken)

	authService := service.NewAuthService(logger, mockUserRepo, mockToken)

	token := "valid-token"
	userID := 1

	// Mock the VerifyToken to return a userID
	mockToken.On("VerifyToken", token).Return(userID, nil)

	verifiedUserID, err := authService.VerifyJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, userID, verifiedUserID)

	mockToken.AssertExpectations(t)
}

func TestAuthService_VerifyJWT_InvalidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	mockUserRepo := new(mockRepo.MockUserRepository)
	mockToken := new(mockService.MockToken)

	authService := service.NewAuthService(logger, mockUserRepo, mockToken)

	token := "invalid-token"

	// Mock the VerifyToken to return an error
	mockToken.On("VerifyToken", token).Return(-1, errors.New("invalid token"))

	verifiedUserID, err := authService.VerifyJWT(token)

	assert.Error(t, err)
	assert.Equal(t, -1, verifiedUserID)

	mockToken.AssertExpectations(t)
}
