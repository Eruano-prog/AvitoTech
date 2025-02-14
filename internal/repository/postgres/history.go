// Package postgres
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

func (h History) InsertOperation(operation entity.Operation) (*entity.Operation, error) {
	q, err := h.db.Prepare(`
	INSERT INTO history (sender_name, receiver_name, amount)
	VALUES ($1, $2, $3)
	RETURNING id, sender_name, receiver_name, amount
`)
	if err != nil {
		h.l.Error("Failed to prepare query", zap.Error(err))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			h.l.Error("Failed to close query", zap.Error(err))
		}
	}(q)

	var op entity.Operation
	err = q.QueryRow(operation.FromUser, operation.ToUser, operation.Amount).Scan(&op.ID, &op.FromUser, &op.ToUser, &op.Amount)
	if err != nil {
		h.l.Error("Failed to insert history", zap.Error(err))
		return nil, err
	}

	return &op, nil
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
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			h.l.Error("Failed to close query", zap.Error(err))
		}
	}(q)

	rows, err := q.Query(name)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			h.l.Error("Failed to close rows query", zap.Error(err))
		}
	}(rows)

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
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			h.l.Error("Failed to close query", zap.Error(err))
		}
	}(q)

	rows, err := q.Query(name)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			h.l.Error("Failed to close rows query", zap.Error(err))
		}
	}(rows)

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

func (h History) DeleteOperation(id int) error {
	q, err := h.db.Prepare(`
	DELETE FROM history
	WHERE id = $1
`)
	if err != nil {
		h.l.Error("Failed to prepare query", zap.Error(err))
		return err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			h.l.Error("Failed to close query", zap.Error(err))
		}
	}(q)

	_, err = q.Exec(id)
	if err != nil {
		h.l.Error("Failed to delete history", zap.Error(err))
		return err
	}
	return nil
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
