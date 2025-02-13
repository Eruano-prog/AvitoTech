package postgres

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"database/sql"
	"go.uber.org/zap"
)

type History struct {
	l  *zap.Logger
	db *sql.DB
}

func (h History) InsertOperation(operation entity.Operation) error {
	q, err := h.db.Prepare(`
	INSERT INTO history (sender_name, receiver_name, amount)
	VALUES ($1, $2, $3)
`)
	if err != nil {
		h.l.Error("Failed to prepare query", zap.Error(err))
		return err
	}
	defer q.Close()

	_, err = q.Exec(operation.FromUser, operation.ToUser, operation.Amount)
	if err != nil {
		h.l.Error("Failed to insert history", zap.Error(err))
		return err
	}
	return nil
}

func (h History) GetSentByUser(name string) ([]entity.Operation, error) {
	q, err := h.db.Prepare(`
	SELECT sender_name, receiver_name, amount
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
		err = rows.Scan(&operation.FromUser, &operation.ToUser, &operation.Amount)
		if err != nil {
			h.l.Debug("Error scanning rows", zap.Error(err))
			return nil, err
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

func (h History) GetReceivedByUser(name string) ([]entity.Operation, error) {
	q, err := h.db.Prepare(`
	SELECT sender_name, receiver_name, amount
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
		err = rows.Scan(&operation.FromUser, &operation.ToUser, &operation.Amount)
		if err != nil {
			h.l.Debug("Error scanning rows", zap.Error(err))
			return nil, err
		}
		operations = append(operations, operation)
	}
	return operations, nil
}

func NewHistoryRepository(
	l *zap.Logger,
	db *sql.DB,
) repository.HistoryRepository {
	return &History{
		l:  l,
		db: db,
	}
}
