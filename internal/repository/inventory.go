package repository

import (
	"database/sql"
	"fmt"
	"go.uber.org/zap"
)

type InventoryRepository struct {
	l  *zap.Logger
	db *sql.DB
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
	defer q.Close()

	rows, err := q.Query(userID)
	if err != nil {
		i.l.Error("failed to query", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

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

func NewInventoryRepository(
	l *zap.Logger,
	pgAddress string,
	pgUser string,
	pgPassword string,
	pgDatabase string,
) (*InventoryRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s/%s", pgUser, pgPassword, pgAddress, pgDatabase)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Fatal("failed to connect to database", zap.String("dsn", dsn), zap.Error(err))
		return nil, err
	}
	return &InventoryRepository{
		l:  l,
		db: db,
	}, nil
}
