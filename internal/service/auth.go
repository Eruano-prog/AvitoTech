package service

import "errors"

var (
	UnauthorizedError = errors.New("unauthorized")
)

type AuthService struct {
}

// Authenticate returns token associated with user
func (a AuthService) Authenticate(username, password string) (string, error) {
	panic("implement me")
}

// VerifyJWT returns userId if succeeded. If not returns err
func (a AuthService) VerifyJWT(token string) (int, error) {
	panic("implement me")
}
