// Package service
package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUnauthorized     = errors.New("unauthorized")
	ErrUserAlreadyExist = errors.New("user already exist")
)

type AuthService struct {
	l              *zap.Logger
	jwtService     Token
	userRepository repository.UserRepository
}

func (a AuthService) createUser(username, password string) (*entity.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.l.Error("failed to hash password", zap.Error(err))
		return nil, err
	}

	user := &entity.User{
		Username: username,
		Password: string(hashedPassword),
		Balance:  1000,
	}

	user, err = a.userRepository.InsertUser(user)
	if err != nil {
		a.l.Error("failed to insert user", zap.Error(err))
		return nil, err
	}

	return user, nil
}

// Authenticate returns token associated with user
func (a AuthService) Authenticate(username, password string) (string, error) {
	user, err := a.userRepository.FindUserByUsername(username)

	if errors.Is(err, repository.ErrorUserNotFound) {
		user, err = a.createUser(username, password)
		if err != nil {
			return "", err
		}

		token, err := a.jwtService.GenerateToken(user.ID)
		if err != nil {
			a.l.Error("failed to generate token", zap.Error(err))
			return "", err
		}

		return token, nil
	}
	if err != nil {
		a.l.Error("failed to find user by username", zap.Error(err))
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		a.l.Debug("failed to compare password", zap.Error(err))
		return "", ErrUnauthorized
	}

	token, err := a.jwtService.GenerateToken(user.ID)
	if err != nil {
		a.l.Error("failed to generate token", zap.Error(err))
		return "", err
	}
	return token, nil
}

// VerifyJWT returns userId if succeeded. If not returns err
func (a AuthService) VerifyJWT(token string) (int, error) {
	return a.jwtService.VerifyToken(token)
}

func NewAuthService(
	l *zap.Logger,
	u repository.UserRepository,
	j Token,
) Auth {
	return &AuthService{
		l:              l,
		userRepository: u,
		jwtService:     j,
	}
}
