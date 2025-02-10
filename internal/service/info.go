package service

import (
	"AvitoTech/internal/entity"
	"go.uber.org/zap"
)

type InfoService struct {
	l *zap.Logger
}

func (i InfoService) GetInfo(userID int) (entity.AccountInfo, error) {
	panic("implement me")
}

func NewInfoService(l *zap.Logger) *InfoService {
	return &InfoService{l: l}
}
