package repository

import (
	"AvitoTech/internal/entity"
	"errors"
)

var ErrorUserNotFound = errors.New("user not found")

type HistoryRepository interface {
	InsertOperation(operation entity.Operation) error
	GetSentByUser(name string) ([]entity.Operation, error)
	GetReceivedByUser(name string) ([]entity.Operation, error)
}

type InventoryRepository interface {
	InsertItem(owner int, item string) error
	GetUsersInventory(userID int) (map[string]int, error)
}
type UserRepository interface {
	InsertUser(user *entity.User) (*entity.User, error)
	FindUserByUsername(username string) (*entity.User, error)
	FindUserById(id int) (*entity.User, error)
	TransferMoney(userFrom int, userTo int, amount int) error
	WithdrawMoney(user int, amount int) error
}
