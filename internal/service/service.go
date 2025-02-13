package service

import "AvitoTech/internal/entity"

type Auth interface {
	createUser(username, password string) (*entity.User, error)
	Authenticate(username, password string) (string, error)
	VerifyJWT(token string) (int, error)
}
type Token interface {
	GenerateToken(userID int) (string, error)
	VerifyToken(tokenString string) (int, error)
}
type Info interface {
	GetInfo(userID int) (*entity.AccountInfo, error)
}
type Coin interface {
	SendCoin(fromUser int, toUser string, amount int) error
	BuyItem(id int, item string) error
}
