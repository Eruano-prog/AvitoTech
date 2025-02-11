package repository

import (
	"database/sql"
	"go.uber.org/zap"
)

type HistoryRepository struct {
	l  *zap.Logger
	db *sql.DB
}

func NewHistoryRepository(
	l *zap.Logger,
	pgAddress string,
	pgUser string,
	pgPassword string,
	pgDatabase string,
) (*UserDb, error) {
) *HistoryRepository {

}
