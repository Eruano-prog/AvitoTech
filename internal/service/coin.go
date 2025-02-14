package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"errors"
	"go.uber.org/zap"
)

type CoinService struct {
	l *zap.Logger

	userRepo      repository.UserRepository
	inventoryRepo repository.InventoryRepository
	historyRepo   repository.HistoryRepository
}

func (c CoinService) SendCoin(fromUser int, toUser string, amount int) error {
	sender, err := c.userRepo.FindUserById(fromUser)
	if err != nil {
		c.l.Debug("fromUser not found", zap.Error(err))
		return err
	}

	receiver, err := c.userRepo.FindUserByUsername(toUser)
	if err != nil {
		c.l.Debug("toUser not found", zap.Error(err))
		return err
	}

	err = c.userRepo.TransferMoney(fromUser, receiver.Id, amount)
	if err != nil {
		c.l.Debug("failed to transfer money", zap.Error(err))
		return err
	}

	_, err = c.historyRepo.InsertOperation(entity.Operation{
		FromUser: sender.Username,
		ToUser:   receiver.Username,
		Amount:   amount,
	})
	if err != nil {
		c.l.Debug("failed to insert history", zap.Error(err))
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

	_, err = c.inventoryRepo.InsertItem(id, item)
	if err != nil {
		c.l.Error("failed to insert item", zap.Error(err))
		return err
	}

	return nil
}

func NewCoinService(
	l *zap.Logger,
	u repository.UserRepository,
	i repository.InventoryRepository,
	h repository.HistoryRepository,
) Coin {
	return &CoinService{
		l:             l,
		userRepo:      u,
		inventoryRepo: i,
		historyRepo:   h,
	}
}
