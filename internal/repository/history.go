package repository

import (
	"AvitoTech/internal/entity"
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type HistoryRepository struct {
	l  *zap.Logger
	db *sql.DB
}

func (r HistoryRepository) GetSentByUser(name string) ([]entity.Operation, error) {
	q, err := r.db.Prepare(`
	SELECT receiver_name, amount
	FROM history
	WHERE sender_name = $1
`)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	rows, err := q.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []entity.Operation
	for rows.Next() {
		var operation entity.Operation
		err = rows.Scan(&operation.User, &operation.Amount)
		if err != nil {
			r.l.Debug("Error scanning rows", zap.Error(err))
			return nil, err
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

func (r HistoryRepository) GetReceivedByUser(name string) ([]entity.Operation, error) {
	q, err := r.db.Prepare(`
	SELECT sender_name, amount
	FROM history
	WHERE receiver_name = $1
`)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	rows, err := q.Query(name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var operations []entity.Operation
	for rows.Next() {
		var operation entity.Operation
		err = rows.Scan(&operation.User, &operation.Amount)
		if err != nil {
			r.l.Debug("Error scanning rows", zap.Error(err))
			return nil, err
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

func NewHistoryRepository(
	l *zap.Logger,
	pgAddress string,
	pgUser string,
	pgPassword string,
	pgDatabase string,
) (*HistoryRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPassword, pgAddress, pgDatabase)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Fatal("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}

	return &HistoryRepository{
		l:  l,
		db: db,
	}, nil
}
