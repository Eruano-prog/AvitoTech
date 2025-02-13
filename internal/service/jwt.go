package service

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
)

type JWTService struct {
	l      *zap.Logger
	secret []byte
}

func (s JWTService) GenerateToken(userID int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userID,
	})

	return token.SignedString(s.secret)
}

// VerifyToken validates token and return userID from claims
func (s JWTService) VerifyToken(tokenString string) (int, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.MapClaims{}, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		s.l.Error("Error parsing token", zap.Error(err))
		return -1, err
	}

	if claims, ok := token.Claims.(*jwt.MapClaims); ok && token.Valid {
		id, ok := (*claims)["userID"].(float64)
		if !ok {
			s.l.Error("Invalid userID type")
			return -1, errors.New("invalid userID type")
		}
		return int(id), nil
	}

	s.l.Error("Error verifying token", zap.Error(err))
	return -1, errors.New("error verifying token")
}

func NewJWTService(l *zap.Logger, s string) Token {
	return &JWTService{
		l:      l,
		secret: []byte(s),
	}
}
