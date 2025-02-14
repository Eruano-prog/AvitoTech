// Package repository
package repository

import (
	"AvitoTech/internal/entity"
	"errors"
)

var ErrorUserNotFound = errors.New("user not found")

type HistoryRepository interface {
	InsertOperation(operation entity.Operation) (*entity.Operation, error)
	GetSentByUser(name string) ([]entity.Operation, error)
	GetReceivedByUser(name string) ([]entity.Operation, error)
	DeleteOperation(id int) error
}

type InventoryRepository interface {
	InsertItem(owner int, item string) (*entity.Item, error)
	GetUsersInventory(userID int) (map[string]int, error)
	DeleteItem(id int) error
}
type UserRepository interface {
	InsertUser(user *entity.User) (*entity.User, error)
	FindUserByUsername(username string) (*entity.User, error)
	FindUserByID(id int) (*entity.User, error)
	TransferMoney(userFrom int, userTo int, amount int) error
	WithdrawMoney(user int, amount int) error
}
