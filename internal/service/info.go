package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"go.uber.org/zap"
)

type InfoService struct {
	l *zap.Logger

	userRepo      repository.UserRepository
	historyRepo   repository.HistoryRepository
	inventoryRepo repository.InventoryRepository
}

func (i InfoService) GetInfo(userID int) (*entity.AccountInfo, error) {
	user, err := i.userRepo.FindUserByID(userID)
	if err != nil {
		i.l.Debug("user not found", zap.Error(err))
		return nil, err
	}

	sent, err := i.historyRepo.GetSentByUser(user.Username)
	if err != nil {
		i.l.Debug("history not found", zap.Error(err))
	}
	received, err := i.historyRepo.GetReceivedByUser(user.Username)
	if err != nil {
		i.l.Debug("history not found", zap.Error(err))
	}

	inventory, err := i.inventoryRepo.GetUsersInventory(user.ID)
	if err != nil {
		i.l.Debug("inventory not found", zap.Error(err))
	}

	return &entity.AccountInfo{
		Coins:     user.Balance,
		Sent:      sent,
		Received:  received,
		Inventory: inventory,
	}, nil
}

func NewInfoService(
	l *zap.Logger,
	u repository.UserRepository,
	h repository.HistoryRepository,
	i repository.InventoryRepository,
) Info {
	return &InfoService{
		l:             l,
		userRepo:      u,
		historyRepo:   h,
		inventoryRepo: i,
	}
}
