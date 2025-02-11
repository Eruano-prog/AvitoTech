package service

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"go.uber.org/zap"
)

type InfoService struct {
	l *zap.Logger

	userRepo *repository.UserDb
}

func (i InfoService) GetInfo(userID int) (entity.AccountInfo, error) {
	user, err := i.userRepo.FindUserById(userID)
	if err != nil {
		i.l.Debug("user not found", zap.Error(err))
		return entity.AccountInfo{}, err
	}

	se

	result := entity.AccountInfo{
		Coins: user.Balance,
	}
}

func NewInfoService(
	l *zap.Logger,
	u *repository.UserDb,
) *InfoService {
	return &InfoService{
		l:        l,
		userRepo: u,
	}
}
