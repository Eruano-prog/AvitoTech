package service

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestJWTService_GenerateToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	secret := "my-secret-key"
	s := NewJWTService(logger, secret)

	userID := 123
	token, err := s.GenerateToken(userID)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	parsedUserID, err := s.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTService_VerifyToken_ValidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	secret := "my-secret-key"
	s := NewJWTService(logger, secret)

	userID := 123
	token, err := s.GenerateToken(userID)
	assert.NoError(t, err)

	parsedUserID, err := s.VerifyToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, parsedUserID)
}

func TestJWTService_VerifyToken_InvalidToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	secret := "my-secret-key"
	s := NewJWTService(logger, secret)

	invalidToken := "invalid.token.here"
	parsedUserID, err := s.VerifyToken(invalidToken)

	assert.Error(t, err)
	assert.Equal(t, -1, parsedUserID)
}

func TestJWTService_VerifyToken_InvalidUserIDType(t *testing.T) {
	logger, _ := zap.NewProduction()
	secret := "my-secret-key"
	s := NewJWTService(logger, secret)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": "not-an-int",
	})
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	parsedUserID, err := s.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, -1, parsedUserID)
	assert.Contains(t, err.Error(), "invalid userID type")
}

func TestJWTService_VerifyToken_ExpiredToken(t *testing.T) {
	logger, _ := zap.NewProduction()
	secret := "my-secret-key"
	s := NewJWTService(logger, secret)

	claims := jwt.MapClaims{
		"userID": 123,
		"exp":    jwt.NewNumericDate(time.Now().Add(-time.Hour)),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)

	parsedUserID, err := s.VerifyToken(tokenString)
	assert.Error(t, err)
	assert.Equal(t, -1, parsedUserID)
	assert.Contains(t, err.Error(), "token is expired")
}
