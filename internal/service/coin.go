package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"errors"
	"go.uber.org/zap"
)

type CoinService struct {
	l *zap.Logger

	userRepo      *repository.UserRepository
	inventoryRepo *repository.InventoryRepository
}

func (c CoinService) SendCoin(fromUser int, toUser string, amount int) error {
	receiver, err := c.userRepo.FindUserByUsername(toUser)
	if err != nil {
		c.l.Debug("toUser bot found")
		return err
	}

	err = c.userRepo.TransferMoney(fromUser, receiver.Id, amount)
	if err != nil {
		c.l.Debug("failed to transfer money", zap.Error(err))
		return err
	}

	return nil
}

func (c CoinService) BuyItem(id int, item string) error {
	cost, exist := entity.Items[item]
	if !exist {
		return errors.New("item not found")
	}

	err := c.userRepo.WithdrawMoney(id, cost)
	if err != nil {
		c.l.Error("failed to withdrawMoney", zap.Error(err))
		return err
	}

	err = c.inventoryRepo.InsertItem(id, item)
	if err != nil {
		c.l.Error("failed to insert item", zap.Error(err))
		return err
	}

	return nil
}

func NewCoinService(
	l *zap.Logger,
	u *repository.UserRepository,
	i *repository.InventoryRepository,
) *CoinService {
	return &CoinService{
		l:             l,
		userRepo:      u,
		inventoryRepo: i,
	}
}
