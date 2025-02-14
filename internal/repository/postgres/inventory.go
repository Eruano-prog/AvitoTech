package postgres

import (
	"AvitoTech/internal/entity"
	"AvitoTech/internal/repository"
	"database/sql"
	"go.uber.org/zap"
)

type InventoryRepository struct {
	l  *zap.Logger
	db *sql.DB
}

func (i InventoryRepository) InsertItem(owner int, itemTitle string) (*entity.Item, error) {
	q, err := i.db.Prepare(`
	INSERT INTO inventory (owner_id, item)
	VALUES ($1, $2)
	RETURNING id, owner_id, item
`)
	if err != nil {
		i.l.Error("failed to prepare inventory query", zap.Error(err))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			i.l.Error("failed to close inventory query", zap.Error(err))
		}
	}(q)

	var item entity.Item
	err = q.QueryRow(owner, itemTitle).Scan(&item.ID, &item.OwnerID, &item.Title)
	if err != nil {
		i.l.Error("failed to insert item", zap.Error(err))
		return nil, err
	}

	return &item, nil
}

func (i InventoryRepository) GetUsersInventory(userID int) (map[string]int, error) {
	q, err := i.db.Prepare(`
	SELECT item, count(item)
	FROM inventory
	WHERE owner_id = $1
	GROUP BY item
`)
	if err != nil {
		i.l.Error("failed to prepare query", zap.Error(err))
		return nil, err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			i.l.Error("failed to close inventory query", zap.Error(err))
		}
	}(q)

	rows, err := q.Query(userID)
	if err != nil {
		i.l.Error("failed to query", zap.Error(err))
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err = rows.Close()
		if err != nil {
			i.l.Error("failed to close inventory query", zap.Error(err))
		}
	}(rows)

	var result = make(map[string]int)
	for rows.Next() {
		var count int
		var item string
		err = rows.Scan(&item, &count)
		if err != nil {
			i.l.Error("failed to scan", zap.Error(err))
			return nil, err
		}

		result[item] = count
	}

	return result, nil
}

func (i InventoryRepository) DeleteItem(id int) error {
	q, err := i.db.Prepare(`
	DELETE FROM inventory
	WHERE id = $1
`)
	if err != nil {
		i.l.Error("failed to prepare delete query", zap.Error(err))
		return err
	}
	defer func(q *sql.Stmt) {
		err = q.Close()
		if err != nil {
			i.l.Error("failed to close inventory query", zap.Error(err))
		}
	}(q)

	_, err = q.Exec(id)
	if err != nil {
		i.l.Error("failed to delete item", zap.Error(err))
		return err
	}
	return nil
}

func NewInventoryRepository(
	l *zap.Logger,
	db *sql.DB,
) repository.InventoryRepository {
	return &InventoryRepository{
		l:  l,
		db: db,
	}
}
