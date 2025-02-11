package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"errors"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

var (
	UnauthorizedError     = errors.New("unauthorized")
	UserAlreadyExistError = errors.New("user already exist")
)

type AuthService struct {
	l              *zap.Logger
	jwtService     *JWTService
	userRepository *repository.UserDb
}

// Authenticate returns token associated with user
// TODO: DEBUG
func (a AuthService) Authenticate(username, password string) (string, error) {
	_, err := a.userRepository.FindUserByUsername(username)
	if err == nil {
		return "", UserAlreadyExistError
	}

	encPasswd, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		a.l.Error("failed to hash password", zap.Error(err))
		return "", err
	}

	user := &entity.User{
		Username: username,
		Password: string(encPasswd),
		Balance:  0,
	}

	user, err = a.userRepository.InsertUser(user)
	if err != nil {
		a.l.Error("failed to insert user", zap.Error(err))
		return "", err
	}

	token, err := a.jwtService.GenerateToken(user.Id)
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
	userRepository *repository.UserDb,
) *AuthService {
	return &AuthService{
		l:              l,
		userRepository: userRepository,
	}
}
